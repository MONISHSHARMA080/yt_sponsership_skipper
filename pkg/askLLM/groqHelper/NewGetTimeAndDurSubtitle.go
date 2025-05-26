package askllm

import "strings"

func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming) {
	newFullCaptions := SanatizeStrignsForSearching(full_captions)
	newsponserSubtitleFromLLM := SanatizeStrignsForSearching(sponserSubtitleFromLLM)
	full_captions = &newFullCaptions
	sponserSubtitleFromLLM = &newsponserSubtitleFromLLM
}

// this will write to the location of the string and this is not a pass by copy,
// this is done to achive speed and reduce the cost, only use this if you are sure that you will not need the ptr string again
// takes the string and cleans it by removign any more than one " "(spaces) and "\n" with " " and also make it in the lower cap
func SanatizeStrignsForSearching(s *string) string {
	newStr := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(*s), "  ", " "), "\n", " ")
	return strings.Join(strings.Fields(newStr), " ")
}

// this one will not use a ptr
// takes the string and cleans it by removign any more than one " "(spaces) and "\n" with " " and also make it in the lower cap
func SanatizeStrignsForSearchingWithoutPtr(s string) string {
	newStr := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(s), "  ", " "), "\n", " ")
	return strings.Join(strings.Fields(newStr), " ")
}
