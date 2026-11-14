package main

import "testing"

func TestAddPatient(t *testing.T) {
	var p []Patient
	patient1 := Patient{ID: 1, Name: "Alice", Age: 30}
	p = addPatient(p, patient1)
	if len(p) != 1 {
		t.Fatalf("Expected 1 patient, got %d", len(p))
	}
	if p[0] != patient1 {
		t.Errorf("Expected patient %+v, got %+v", patient1, p[0])
	}
}

func TestFindPatientByID(t *testing.T) {
	patients := []Patient{
		{ID: 1, Name: "Alice", Age: 20},
		{ID: 2, Name: "Bob", Age: 25},
	}
	foundPatient, err := findPatientByID(patients, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if foundPatient == nil {
		t.Fatalf("Expected patient, got nil")
	}
	if foundPatient.Name != "Bob" {
		t.Errorf("Expected name 'Bob', got '%s'", foundPatient.Name)
	}
	// TODO: 找到 ID=2 的患者
	// 验证 err == nil
	// 验证 Name == "Bob"
}

func TestPatientNotFound(t *testing.T) {
	patients := []Patient{
		{ID: 1, Name: "Alice", Age: 20},
	}
	foundPatient, err := findPatientByID(patients, 999)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if foundPatient != nil {
		t.Errorf("Expected nil, got %v", foundPatient)
	}
	// TODO: 查询 ID=999
	// 验证 err != nil
	// 验证返回的指针 == nil
}

func TestPatientMutation(t *testing.T) {
	patients := []Patient{
		{ID: 1, Name: "Alice", Age: 20},
	}

	foundPatient, err := findPatientByID(patients, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if foundPatient == nil {
		t.Fatalf("Expected patient, got nil")
	}
	foundPatient.Name = "Charlie"
	if patients[0].Name != "Charlie" {
		t.Errorf("Expected name 'Charlie', got '%s'", patients[0].Name)
	}
	//range返回副本地址，但find中我们使用了i来取得原始切片的地址，所以修改foundPatient.Name会影响patients[0].Name
	// TODO: 找到 ID=1
	// 将找到的患者 Name 修改为 "Charlie"
	// 检查 patients[0].Name 是否同步改变
}
