package main

func main() {

	reception := &Reception{}
	doctor := &Doctor{}
	medical := &Medical{}
	cashier := &Cashier{}

	reception.setNext(doctor)
	doctor.setNext(medical)
	medical.setNext(cashier)

	patient := &Patient{name: "abc"}
	//Patient visiting
	reception.execute(patient)
}
