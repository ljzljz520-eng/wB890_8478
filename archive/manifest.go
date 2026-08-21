package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"memorialstation/domain"
	"sort"
)

type Manifest struct {
	BatchID   string
	Entries   int
	Digest    string
	RecordIDs []string
}

func (s *Service) BuildManifest(batchID string) (Manifest, error) {
	entries, err := s.store.ListArchives(batchID)
	if err != nil {
		return Manifest{}, err
	}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.RecordID)
	}
	sort.Strings(ids)
	digest := sha256.Sum256([]byte(join(ids)))
	return Manifest{BatchID: batchID, Entries: len(entries), Digest: hex.EncodeToString(digest[:]), RecordIDs: ids}, nil
}
func join(values []string) string {
	output := ""
	for _, value := range values {
		output += value + "\n"
	}
	return output
}
func (s *Service) VerifyManifest(manifest Manifest) (bool, error) {
	current, err := s.BuildManifest(manifest.BatchID)
	if err != nil {
		return false, err
	}
	if current.Entries != manifest.Entries || current.Digest != manifest.Digest {
		return false, nil
	}
	return true, nil
}
func (s *Service) VisibleSnapshot(batchID string) ([]string, error) {
	records, err := s.ReadBack(batchID, "manifest")
	if err != nil {
		return nil, err
	}
	snapshots := make([]string, 0, len(records))
	for _, record := range records {
		if record.Visibility == domain.VisibilityPublic {
			snapshots = append(snapshots, record.StudentName+":"+record.Message)
		}
	}
	sort.Strings(snapshots)
	return snapshots, nil
}
