package rna

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Diff needs a DiffConfig file to perform.
type Diff struct {
	Config DiffConfig
}

// Configure is a method used before generating pbs files.
func (diff *Diff) Configure() error {
	config := diff.Config
	if config.ProjectPath != "" {
		return nil
	}

	return config.Parse()
}

func genJoinedTophatBamPathes(a []string, projectPath string) string {
	for i := range a {
		a[i] = filepath.Join(projectPath, a[i], "accepted_hits.bam")
	}

	return strings.Join(a, ",")
}

// GenPbsFiles is a method of rna.Diff, which as the name,
// generates several pbs files to perform `cuffdiff` action.
func (diff *Diff) GenPbsFiles(user string) error {
	if err := diff.Configure(); err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	dc := diff.Config
	for _, comp := range dc.Comparisons {
		// Get `Group` by names in `Comparison`s
		caseGrp, caseErr := dc.GetGroupByName(comp.Case)
		if caseErr != nil {
			return caseErr
		}
		ctrlGrp, ctrlErr := dc.GetGroupByName(comp.Control)
		if ctrlErr != nil {
			return ctrlErr
		}

		// Generate diff output path
		outputDir := fmt.Sprintf("%s_%s_DEG", caseGrp.Name, ctrlGrp.Name)
		outputPath, absErr1 := filepath.Abs(filepath.Join(dc.ProjectPath, outputDir))
		if absErr1 != nil {
			return absErr1
		}

		// Generate pbs file path
		pbsFileName := fmt.Sprintf("batch_DEG_%s_%s.pbs", caseGrp.Name, ctrlGrp.Name)
		pbsFileDir := filepath.Join(dc.ProjectPath, "log")
		if err := os.MkdirAll(pbsFileDir, 0774); err != nil {
			return err
		}

		pbsFilePath, absErr2 := filepath.Abs(filepath.Join(pbsFileDir, pbsFileName))
		if absErr2 != nil {
			return absErr2
		}

		pbsFile, pbsErr := os.Create(pbsFilePath)
		defer pbsFile.Close()
		if pbsErr != nil {
			return pbsErr
		}

		pbsFile.WriteString(fmt.Sprintf("#PBS -N DEG_%s_%s\n", caseGrp.Name, ctrlGrp.Name))
		pbsFile.WriteString(fmt.Sprintf("#PBS -o %s/%s/DEG_%s_%s.out\n",
			dc.ProjectPath, "log", caseGrp.Name, ctrlGrp.Name))
		pbsFile.WriteString(fmt.Sprintf("#PBS -e %s/%s/DEG_%s_%s.err\n",
			dc.ProjectPath, "log", caseGrp.Name, ctrlGrp.Name))
		pbsFile.WriteString(fmt.Sprintf("#PBS -l nodes=1:ppn=1\n"))
		pbsFile.WriteString(fmt.Sprintf("#PBS -q high\n"))
		pbsFile.WriteString(fmt.Sprintf("#PBS -r y\n"))
		pbsFile.WriteString(fmt.Sprintf("#PBS -u %s\n\n", user))
		pbsFile.WriteString(fmt.Sprintf("%s --no-update-check -o %s -b %s -u %s %s %s",
			dc.CuffdiffExec, outputPath, dc.GenomeBtwIdx, dc.TranxBtwIdx,
			genJoinedTophatBamPathes(caseGrp.Samples, dc.ProjectPath),
			genJoinedTophatBamPathes(ctrlGrp.Samples, dc.ProjectPath)))

		fmt.Println("qsub", pbsFilePath)
	}

	return nil
}
