package askllm

import (
	"strings"

	"go.uber.org/zap"
)

// func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming, logger *zap.Logger) {
func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, logger *zap.Logger) {
	newFullCaptions := SanatizeStrignsForSearching(full_captions)
	newsponserSubtitleFromLLM := SanatizeStrignsForSearching(sponserSubtitleFromLLM)
	full_captions = &newFullCaptions
	sponserSubtitleFromLLM = &newsponserSubtitleFromLLM

	// now I need to get the start and end time/index(both) from the transcript
	// I know that both of the transcripts.subtitles and the full caption and the sponserSubtitleFromLLM is sanitized
	// the transcripts do not have the " " space in the end but the full_captions and the sponserSubtitleFromLLM captions does so we will need to take it into account
	// NOTE: the strings are utf-8 encoded
	//
	//we are going to get the subtitles using the lenght technique:--> first we will get the sponsership start index in the full captions and  the we will get the end index by
	// doing start index plus the string lenght (make sure to get the last " " ),
	sponsershipStartIndex := strings.Index(*full_captions, *sponserSubtitleFromLLM)
	logger.Info("got the sponsershipStartIndex ", zap.Int("the sponsershipStartIndex is ", sponsershipStartIndex))
	if sponsershipStartIndex == -1 {
		logger.Info("failed getting the logs, as the index() returned -1", zap.Skip())
		// channel response
		// responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{0, 0, fmt.Errorf("can't find the subtitles the result for the start index is -1")}
		return
	}
	sponserShipEndIndex := sponsershipStartIndex + len(*sponserSubtitleFromLLM)
	println("the sponsershipStartIndex is:", sponsershipStartIndex, " and the end index:", sponserShipEndIndex, " is the sponserSubtitleFromLLM == full_captions[sponsershipStartIndex : sponserShipEndIndex ]",
		*sponserSubtitleFromLLM == (*full_captions)[sponsershipStartIndex:sponserShipEndIndex])

	// used for collecting the index of the start and the end
	type sponsershipPositionIndex struct {
		indexAt int // contains the subtitle array's index
		found   bool
	}
	sponserShipStartIndexForSubtitlesArray := sponsershipPositionIndex{found: false, indexAt: 0}
	sponserShipEndIndexForSubtitlesArray := sponsershipPositionIndex{found: false, indexAt: 0}

	strInSubtitlesLenCounter := 0
	for i, subtitle := range transcripts.Subtitles {
		strInSubtitlesLenCounter += len(subtitle.Text + " ")
		// if we are after the sponsor ship start index and it is not assigned before
		if !sponserShipStartIndexForSubtitlesArray.found && strInSubtitlesLenCounter >= sponsershipStartIndex {
			sponserShipStartIndexForSubtitlesArray.found = true
			sponserShipStartIndexForSubtitlesArray.indexAt = strInSubtitlesLenCounter
			println("we have reached the sponserShipStartIndex and it is at the index:", i, " and the text in this index is:=>", subtitle.Text)
		}
		if !sponserShipEndIndexForSubtitlesArray.found && strInSubtitlesLenCounter >= sponserShipEndIndex {
			println("we have reached the sponserShipEndIndex and it is at the index:", i, " and the text in this index is:=>", subtitle.Text)
			sponserShipEndIndexForSubtitlesArray.found = true
			sponserShipEndIndexForSubtitlesArray.indexAt = sponserShipEndIndex
			break // as we have completed the loop
		}
		// as the transcripts do not have the " " but the full_captions and the sponserSubtitleFromLLM captions does so we will need to take it into account
		// I am not doing it at the start of the loop as I do not want to miss the 0 case too
	}
	// ok now that I seems to be almost right, now I need to get the subtitles form the adjacent indexes
	logger.Info("got the index of the start and the end subtitles, now goting to get the actual subtitles from the adjacent indices")

	//
	//
	//
	// -- in the end just make rough assert that the strings.contains() sponserSubtitleFromLLM and the rought string that you generated to make sure that we are right
	//  for eg if the sponser is form index 1 to 10 I might get the 90 percent accuraccy as there might be some strings there too, but all in all (create a massive string form start to
	//  end for testing/assrting) is mostly right, as that huge generated string will have the sponserSubtitleFromLLM string in it
	// --
	// -----quick assertion testing(see above comment) to make sure tings are alwright --
	// .........To Implement.......
	//
	//
	//
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
