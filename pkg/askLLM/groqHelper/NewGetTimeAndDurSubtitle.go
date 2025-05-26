package askllm

import "strings"

func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming) {
	newFullCaptions := SanatizeStrignsForSearching(full_captions)
	newsponserSubtitleFromLLM := SanatizeStrignsForSearching(sponserSubtitleFromLLM)
	full_captions = &newFullCaptions
	sponserSubtitleFromLLM = &newsponserSubtitleFromLLM

	// now I need to get the start and end time/index(both) from the transcript
	// I know that both of the transcripts.subtitles and the full caption and the sponserSubtitleFromLLM is sanitized
	// NOTE: the strings are utf-8 encoded
	//
	//we are going to get the subtitles using the lenght technique:--> first we will get the sponsership start index in the full captions and  the we will get the end index by
	// doing start index plus the string lenght (make sure to get the last " " , haven't thought this through ), then we just take this out in a var
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
