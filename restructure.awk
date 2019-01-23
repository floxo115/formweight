BEGIN {
  print "type Element struct {"
  print "\tName string"
  print "\tSymbol string"
  print "\tAtomicNumber int"
  print "\tWeight float64"
  print "}"

  print ""
  print "var MapOfEls map[string]*Element"

  print "func init() {"
  print "\tMapOfEls = make(map[string]*Element)"
}

{
  printf "\tMapOfEls[\"%s\"] = &Element{\n", $2 
  printf "\t\tName: \"%s\",\n", $1
  printf "\t\tSymbol: \"%s\",\n", $2
  printf "\t\tAtomicNumber: %s,\n", $3
  printf "\t\tWeight: %s,\n", $4
  print "\t}" 
}

END {
  print "}"
}
