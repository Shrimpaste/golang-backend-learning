package main

import "fmt"

type Patient struct {
	ID   int
	Name string
	Age  int
}

func addPatient(patients []Patient, patient Patient) []Patient {
	return append(patients, patient)
}

func findPatientByID(patients []Patient, id int) (*Patient, error) {
	/*for _, patient := range patients {
		if patient.ID == id {
			return &patient, nil
		}
	}*/
	// range 存的是副本，返回的是副本的地址
	for i := range patients {
		if patients[i].ID == id {
			return &patients[i], nil
		}
	}
	return nil, fmt.Errorf("patient not found")
}

func printPatient(patients []Patient) {
	for _, patient := range patients {
		fmt.Printf("Patient ID: %d, Name: %s, Age: %d\n", patient.ID, patient.Name, patient.Age)
	}
}

func main() {
	var patients []Patient

	patient1 := Patient{ID: 1, Name: "Alan", Age: 30}
	patient2 := Patient{ID: 2, Name: "Bob", Age: 25}
	patient3 := Patient{ID: 3, Name: "Coke", Age: 35}

	patients = addPatient(patients, patient1)
	patients = addPatient(patients, patient2)
	patients = addPatient(patients, patient3)

	printPatient(patients)

	fmt.Printf("Total: %d\n", len(patients))

	foundPatient, err := findPatientByID(patients, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		printPatient([]Patient{*foundPatient})
	}

	foundPatient.Name = "Bob Updated"
	printPatient(patients)
	// 修改了 foundPatient 的 Name 属性，patients 中对应的 Patient 也会被修改，但如果我们在 findPatientByID 中使用了 range 循环的副本方式，那么修改 foundPatient 的 Name 属性不会影响 patients 中的 Patient，因为 foundPatient 是一个副本的地址。
}
