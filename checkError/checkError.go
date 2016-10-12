/*
	Package checkError 
	Handle errors 
*/
package checkError

//Check error
func Check(e error) {
	if e != nil {
		panic(e)
	}	
}