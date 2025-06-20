package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// --- Test Setup with FakeLogReader ---

// FakeLogReader is a test implementation of the LogReader interface.
type FakeLogReader struct {
	// logs maps a pod identifier ("namespace/podName") to its log content.
	logs map[string]string
}

// newFakeLogReader creates an empty FakeLogReader.
func newFakeLogReader() *FakeLogReader {
	return &FakeLogReader{
		logs: make(map[string]string),
	}
}

// addLog adds log content for a specific pod.
func (f *FakeLogReader) addLog(namespace, podName, containerName, content string) {
	key := fmt.Sprintf("%s/%s/%s", namespace, podName, containerName)
	f.logs[key] = content
}

// GetPodLogs retrieves the stored log content for a pod.
func (f *FakeLogReader) GetPodLogs(namespace, podName, containerName string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", namespace, podName, containerName)
	if content, found := f.logs[key]; found {
		return content, nil
	}
	// Return an error if no logs are found for the pod, simulating a real-world scenario.
	return "", fmt.Errorf("logs not found for pod %s", podName)
}

// --- Tests ---

func TestLogGrepper_InteractionWithAPIServer(t *testing.T) {
	// This test now fully emulates the Java test `testInteractionWithAPIServer`.

	// 1. Define pods
	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "test"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container1"}}},
		Status:     corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "container1"}}},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "test"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "container2"}}},
		Status:     corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "container2"}}},
	}

	// 2. Set up the FakeLogReader with realistic log content
	fakeLogReader := newFakeLogReader()
	fakeLogReader.addLog("test", "pod1", "container1", "Initializing xpto\nxpto initialized\nerror writing to xpto")
	fakeLogReader.addLog("test", "pod2", "container2", "Initializing foo\nfoo initialized\nerror writing to foo\nInitializing bar\nbar initialized\nerror writing to bar")

	// 3. Create a standard fake clientset for pod listing
	fakeClientset := fake.NewSimpleClientset(pod1, pod2)

	// 4. Inject the fake clientset and fake log reader into the Grepper
	grepper := &Grepper{
		clientset: fakeClientset,
		logReader: fakeLogReader,
	}

	// 5. Run the test and assert results
	messages := grepper.Grep("test", "pod", "initialized", "POD_AND_CONTAINER")

	// 6. Assertions match the Java test's expectations
	expectedMessages := []Message{
		{PodName: "pod1", ContainerName: "container1", Message: "xpto initialized", LineNumber: 2},
		{PodName: "pod2", ContainerName: "container2", Message: "foo initialized", LineNumber: 2},
		{PodName: "pod2", ContainerName: "container2", Message: "bar initialized", LineNumber: 5},
	}
	assert.ElementsMatch(t, expectedMessages, messages)
}

func TestLogGrepper_ResourceFiltering(t *testing.T) {
	// 1. Define one pod
	appPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "app-pod-1", Namespace: "default"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "c1"}}},
		Status:     corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "c1"}}},
	}

	// 2. Define its logs
	fakeLogReader := newFakeLogReader()
	fakeLogReader.addLog("default", "app-pod-1", "c1", "app log with pattern")

	// 3. Create dependencies
	fakeClientset := fake.NewSimpleClientset(appPod)
	grepper := &Grepper{
		clientset: fakeClientset,
		logReader: fakeLogReader,
	}

	// 4. Run Grep with explicit namespace and resource filter
	messages := grepper.Grep("default", "app", "pattern", "POD_AND_CONTAINER")

	// 5. Assert that the single log message is found
	assert.Len(t, messages, 1)
	assert.Equal(t, "app log with pattern", messages[0].Message)
}

func TestLogGrepper_SortByMessage(t *testing.T) {
	messages := []Message{
		{Message: "zebra message"},
		{Message: "alpha message"},
		{Message: "beta message"},
	}

	grepper := &Grepper{}
	sortedMessages := grepper.sortMessages(messages, "MESSAGE")

	assert.Equal(t, "alpha message", sortedMessages[0].Message)
	assert.Equal(t, "beta message", sortedMessages[1].Message)
	assert.Equal(t, "zebra message", sortedMessages[2].Message)
}

func TestLogGrepper_GetPodLogs_Error(t *testing.T) {
	// The fake reader will return an error because no logs are added for this pod
	fakeLogReader := newFakeLogReader()
	pod1 := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "test"}}
	fakeClientset := fake.NewSimpleClientset(pod1)

	grepper := &Grepper{
		clientset: fakeClientset,
		logReader: fakeLogReader,
	}

	// Since the log reader returns an error, Grep should return no messages
	messages := grepper.Grep("test", "pod1", "any-pattern", "")
	assert.Empty(t, messages)
}

func TestLogGrepper_SearchLogs_EmptyPattern(t *testing.T) {
	grepper := &Grepper{}
	logContent := "line 1\nline 2\nline 3"
	messages := grepper.searchLogs(logContent, "", "pod1", "c1")

	// An empty pattern should match all lines
	assert.Len(t, messages, 3)
	assert.Equal(t, "line 1", messages[0].Message)
	assert.Equal(t, "line 2", messages[1].Message)
	assert.Equal(t, "line 3", messages[2].Message)
}
