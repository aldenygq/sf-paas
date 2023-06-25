package tools


func InterToStringSlice(inters []interface{}) []string {
	var slc []string = make([]string,0)
	for _, v := range inters {
		slc = append(slc,v.(string))
	}
	return slc
}
