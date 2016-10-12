	Author Louiza HANIS
	Date : 01/09/2016

----------------------------------------------------------------------------	
Purpose 

This software is used in metagenomic field. It takes as input a matrix and calculates correlation pairwise between all its vectors, draw histogam of the distribution and generate an RDF file that feeds Bolt database. 

----------------------------------------------------------------------------
Compilation 

go run main.go -i $path/MatrixFileName -o $path/OutputFileName

Input / Ouput parameters 

-i <$path/MatrixFileName>: The file containing the input matrix  

-o <$path/OutputFileName>: The file containing the output correlation matrix  
