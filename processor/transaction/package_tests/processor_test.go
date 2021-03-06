package package_tests

import (
	"regexp"
	"testing"

	"github.com/elastic/apm-server/config"
	"github.com/elastic/apm-server/processor/transaction"
	"github.com/elastic/apm-server/tests"
)

var (
	backendRequestInfo = []tests.RequestInfo{
		{Name: "TestProcessTransactionFull", Path: "../testdata/transaction/payload.json"},
		{Name: "TestProcessTransactionNullValues", Path: "../testdata/transaction/null_values.json"},
		{Name: "TestProcessSystemNull", Path: "../testdata/transaction/system_null.json"},
		{Name: "TestProcessProcessNull", Path: "../testdata/transaction/process_null.json"},
		{Name: "TestProcessTransactionMinimalSpan", Path: "../testdata/transaction/minimal_span.json"},
		{Name: "TestProcessTransactionMinimalService", Path: "../testdata/transaction/minimal_service.json"},
		{Name: "TestProcessTransactionMinimalProcess", Path: "../testdata/transaction/minimal_process.json"},
		{Name: "TestProcessTransactionEmpty", Path: "../testdata/transaction/transaction_empty_values.json"},
		{Name: "TestProcessTransactionAugmentedIP", Path: "../testdata/transaction/augmented_payload_backend.json"},
	}

	backendRequestInfoIgnoreTimestamp = []tests.RequestInfo{
		{Name: "TestProcessTransactionMinimalPayload", Path: "../testdata/transaction/minimal_payload.json"},
	}

	frontendRequestInfo = []tests.RequestInfo{
		{Name: "TestProcessTransactionFrontend", Path: "../testdata/transaction/frontend.json"},
		{Name: "TestProcessTransactionAugmentedMerge", Path: "../testdata/transaction/augmented_payload_frontend.json"},
		{Name: "TestProcessTransactionAugmented", Path: "../testdata/transaction/augmented_payload_frontend_no_context.json"},
	}
)

// ensure all valid documents pass through the whole validation and transformation process
func TestTransactionProcessorOK(t *testing.T) {
	tests.TestProcessRequests(t, transaction.NewProcessor(), config.Config{}, backendRequestInfo, map[string]string{})
}

func TestMinimalTransactionProcessorOK(t *testing.T) {
	tests.TestProcessRequests(t, transaction.NewProcessor(), config.Config{}, backendRequestInfoIgnoreTimestamp, map[string]string{"@timestamp": "-"})
}

func TestProcessorFrontendOK(t *testing.T) {
	conf := config.Config{
		LibraryPattern:      regexp.MustCompile("/test/e2e|~"),
		ExcludeFromGrouping: regexp.MustCompile("^~/test"),
	}
	tests.TestProcessRequests(t, transaction.NewProcessor(), conf, frontendRequestInfo, map[string]string{"@timestamp": "-"})
}

func BenchmarkBackendProcessor(b *testing.B) {
	tests.BenchmarkProcessRequests(b, transaction.NewProcessor(), config.Config{}, backendRequestInfo)
	tests.BenchmarkProcessRequests(b, transaction.NewProcessor(), config.Config{}, backendRequestInfoIgnoreTimestamp)
}

func BenchmarkFrontendProcessor(b *testing.B) {
	conf := config.Config{
		LibraryPattern:      regexp.MustCompile("/test/e2e|~"),
		ExcludeFromGrouping: regexp.MustCompile("^~/test"),
	}
	tests.BenchmarkProcessRequests(b, transaction.NewProcessor(), conf, frontendRequestInfo)
}
