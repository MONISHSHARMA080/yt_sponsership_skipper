package askllm

import (
	"strings"

	"go.uber.org/zap"
)

// used for collecting the index of the start and the end
type sponsershipPositionIndex struct {
	IndexAt int  `json:"indexAt"` // contains the subtitle array's index. 0 indexed btw
	Found   bool `json:"found"`
	// what is the value of the subtitle's lenght at the current iteration
	SubtitleStringCounterAtTheCurrentIteration int
	// what is the value of the subtitle's lenght at the previous iteration, that is the whole index + " " lenght
	SubtitleStringCounterAtThePreviosLoopIteration int
	// what is the value of the subtitle's lenght at the 2nd previous iteration, that is the whole index + " " lenght
	SubtitleStringCounterAtThe2ndPreviousLoopIteration int
}

// func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming, logger *zap.Logger) {
func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, logger *zap.Logger) {
	newFullCaptions := SanatizeStrignsForSearching(full_captions)
	newsponserSubtitleFromLLM := SanatizeStrignsForSearching(sponserSubtitleFromLLM)
	full_captions = &newFullCaptions
	sponserSubtitleFromLLM = &newsponserSubtitleFromLLM
	println("\n\n ---- the full captions in the lower case is---- \n ", *full_captions, "---- \n\n\n")
	println("\n\n ---- the sponsership_subtitles_form_groq in the lower case is --- \n", *sponserSubtitleFromLLM, "---- \n\n\n")

	// now I need to get the start and end time/index(both) from the transcript
	// I know that both of the transcripts.subtitles and the full caption and the sponserSubtitleFromLLM is sanitized
	// the transcripts do not have the " " space in the end but the full_captions and the sponserSubtitleFromLLM captions does so we will need to take it into account
	//  NOTE: the strings are utf-8 encoded
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

	sponserShipStartIndexForSubtitlesArray := sponsershipPositionIndex{Found: false, IndexAt: 0}
	sponserShipEndIndexForSubtitlesArray := sponsershipPositionIndex{Found: false, IndexAt: 0}

	strInSubtitlesLenCounter := 0
	prevSubLenCounter := 0
	secondLastLenCounter := 0
	for i, subtitle := range transcripts.Subtitles {

		secondLastLenCounter = prevSubLenCounter
		prevSubLenCounter = strInSubtitlesLenCounter

		strInSubtitlesLenCounter += len(subtitle.Text + " ")
		// println("the strInSubtitlesLenCounter is:", strInSubtitlesLenCounter, " and prevSubLenCounter:", prevSubLenCounter, " secondLastLenCounter:", secondLastLenCounter)
		// if we are after the sponsor ship start index and it is not assigned before
		if !sponserShipStartIndexForSubtitlesArray.Found && strInSubtitlesLenCounter >= sponsershipStartIndex {
			sponserShipStartIndexForSubtitlesArray.Found = true
			sponserShipStartIndexForSubtitlesArray.IndexAt = i
			sponserShipStartIndexForSubtitlesArray.SubtitleStringCounterAtTheCurrentIteration = strInSubtitlesLenCounter
			sponserShipStartIndexForSubtitlesArray.SubtitleStringCounterAtThePreviosLoopIteration = prevSubLenCounter
			sponserShipStartIndexForSubtitlesArray.SubtitleStringCounterAtThe2ndPreviousLoopIteration = secondLastLenCounter
			println("we have reached the sponserShipStartIndex and it is at the index:", i, " and the text in this index is:=>", subtitle.Text)
		}
		if !sponserShipEndIndexForSubtitlesArray.Found && strInSubtitlesLenCounter >= sponserShipEndIndex {
			println("we have reached the sponserShipEndIndex and it is at the index:", i, " and the text in this index is:=>", subtitle.Text)
			sponserShipEndIndexForSubtitlesArray.Found = true
			sponserShipEndIndexForSubtitlesArray.IndexAt = i
			sponserShipEndIndexForSubtitlesArray.SubtitleStringCounterAtTheCurrentIteration = strInSubtitlesLenCounter
			sponserShipEndIndexForSubtitlesArray.SubtitleStringCounterAtThePreviosLoopIteration = prevSubLenCounter
			sponserShipEndIndexForSubtitlesArray.SubtitleStringCounterAtThe2ndPreviousLoopIteration = secondLastLenCounter
			break // as we have completed the loop
		}
		// as the transcripts do not have the " " but the full_captions and the sponserSubtitleFromLLM captions does so we will need to take it into account
		// I am not doing it at the start of the loop as I do not want to miss the 0 case too
	}

	// if we did not found the sponseShip index is not found then retuern the error
	if !sponserShipStartIndexForSubtitlesArray.Found {
		logger.Error("we were not able to found the sponsership position(start index) in the transcript ", zap.Any("the sponserShipStartIndexForSubtitlesArray is ", sponserShipStartIndexForSubtitlesArray))
		// write to the channel
		return
	}

	if !sponserShipEndIndexForSubtitlesArray.Found {
		logger.Error("we were not able to found the sponsership position(end index) in the transcript ", zap.Any("the sponserShipEndIndexForSubtitlesArray is ", sponserShipEndIndexForSubtitlesArray))
		// write to the channel
		return
	}

	// ok now that I seems to be almost right, now I need to get the subtitles form the adjacent indexes
	logger.Info("got the index of the start and the end subtitles, now goting to get the actual subtitles from the adjacent indices", zap.Any("sponserShipEndIndexForSubtitlesArray is ", sponserShipEndIndexForSubtitlesArray),
		zap.Any("sponserShipStartIndexForSubtitlesArray is ", sponserShipStartIndexForSubtitlesArray),
	)
	preciseStartIndex := getPreciseIndexOfSubtitle(transcripts, &sponserShipStartIndexForSubtitlesArray, sponsershipStartIndex, true)
	println("=============================")
	preciseEndIndex := getPreciseIndexOfSubtitle(transcripts, &sponserShipEndIndexForSubtitlesArray, sponserShipEndIndex, false)

	// -- in the end just make rough assert that the strings.contains() sponserSubtitleFromLLM and the rought string that you generated to make sure that we are right
	//  for eg if the sponser is form index 1 to 10 I might get the 90 percent accuraccy as there might be some strings there too, but all in all (create a massive string form start to
	//  end for testing/assrting) is mostly right, as that huge generated string will have the sponserSubtitleFromLLM string in it
	// --
	// ----------- or -----------
	// 1st iteration:-----
	// we can also do this in the first loop -> make another loop inside it of the range through the words(range over word + " ") there and when we exactly
	// reach the subtitle length (by exact word ) and get the tinings form there
	// ---------------------------
	// 2nd iteration:------
	// or when we found the subtitle in the first loop then we store the subtitle and in the end make a new function where we go back and then
	// in that function look at those adjacent index
	// ---------------------------
	//
	// -----quick assertion testing(see above comment) to make sure tings are alwright --

	println("the precise start index is:", preciseStartIndex, " and the text there is ->", transcripts.Subtitles[preciseStartIndex].Text, " \n\n")
	println("the precise end index is:", preciseEndIndex, " and the text there is ->", transcripts.Subtitles[preciseEndIndex].Text, " \n\n ")
}

// this func is designed to be used when we have gotten the subtitle but it there is a possibility that it could be a error as it might be in the adjacent index
func getPreciseIndexOfSubtitle(transcript *Transcripts, sponserShipIndexInSubtitlesArray *sponsershipPositionIndex, sponserShipAt int, isTheSponserShipPositionForTheStartIndex bool) int {
	// this is a  generic, so it can be used as the first one and the last one
	// ok get the index and form the index-1 to index +1 if we are not able to get it back then we return error
	//
	// the sponserShipAt is the position where the subtitle starts or ends
	// +2 here as we want the range to be i-1 , i , i+1
	println("in the getPreciseIndexOfSubtitle we assert that the sponserShipAt< sponserShipIndexInSubtitlesArray.IndexAt : ", sponserShipAt < sponserShipIndexInSubtitlesArray.IndexAt, "\n")
	sponserShipIndexTracker := sponserShipIndexInSubtitlesArray.SubtitleStringCounterAtThe2ndPreviousLoopIteration
	for i := sponserShipIndexInSubtitlesArray.IndexAt - 1; i < sponserShipIndexInSubtitlesArray.IndexAt+2; i++ {
		// iterate over all the string to reach the precise word in the sponsership (or to the sponserShipAt)
		// now we have a problem do we want to include the space here or not (in the starting), the transcripts do not have space in the begining
		// so we will add the

		subtitleArray := strings.Fields(transcript.Subtitles[i].Text)
		for j, wordWithoutSpace := range subtitleArray {
			// now in the parent function we have the space counted in the previous and the present counter value
			// so at the start here do not include it as it will lead to the double counting of it, and in the  outer loop iteration take the " " in to it as well for the new line and the
			// new index in the subtitle does not have it
			// and to abstract put it above/below the function and then we can just have the if in the inner loop and will all work
			//
			//the problem was that if we have the start index in the first spot in the next subtitle index and the sponserShipIndexTracker is stuck on the index -1
			//so the problem lies that the end index is to the end of the word but the start one is in the startIndex+1 that should solve it
			sponserShipIndexTracker += len(wordWithoutSpace)
			// this is for the current index only
			if sponserShipIndexTracker == sponserShipAt {
				println(" in the case where sponserShipIndexTracker == sponserShipAt and we got the exact position")
				println("-the sponserShipAt(index):->", sponserShipAt, " the current index of words is (sponserShipIndexTracker, track of if we reached sponserShipAt)  ", sponserShipIndexTracker, " and the index found at(rough one) the sponserShipIndexInSubtitlesArray.IndexAt", sponserShipIndexInSubtitlesArray.IndexAt, "word at:", j, " at is ->", wordWithoutSpace, "--   ")
				println(" len(transcript.Subtitles[i].Text) - sponserShipIndexTracker", len(transcript.Subtitles[i].Text)-sponserShipIndexTracker, " \n ")
				// the reason for this if statement is that the  sponserShipIndexTracker == sponserShipAt means we are only one index(in strings) behind the sponserShip start and
				// in the edge case (assume last word in the end of the subtitles index) the sponseShip is in the next index and to be presice we need to skip over
				// but in the case it is for the end function there the end index is the end. These differences is cause of how the len fuc work and how we decided to calculate the
				// start index and the end one(strings.Index gived you the start but the index + full_captions gives you the end word)
				// if isTheSponserShipPositionForTheStartIndex {
				// 	switch lengthLeftInSubArrForMoreString := len(transcript.Subtitles[i].Text) - sponserShipIndexTracker; {
				// 	case lengthLeftInSubArrForMoreString >= 1:
				// 		// return the current index
				// 		return i
				// 		println("we are in the section where there is lengthLeftInSubArrForMoreString >= 1 is true, btw the full subtitle in the row is ->", transcript.Subtitles[i].Text)
				// 	case lengthLeftInSubArrForMoreString < 1:
				// 		println("we are in the section  lengthLeftInSubArrForMoreString<1 and the full subtitles is ->", transcript.Subtitles[i].Text)
				// 		// there is no base case by(mathematically) btw
				// 		if len(transcript.Subtitles) > i+1 {
				// 			return i + 1
				// 		} else {
				// 			return i
				// 		}
				// 		// return
				// 	}
				// } else {
				// 	// we are returning the current index as in the case of the end index this is the last string, unlike the start part where the index is the part before
				// 	// the start of word
				// 	return i
				// }
				switch lengthLeftInSubArrForMoreString := len(transcript.Subtitles[i].Text) - sponserShipIndexTracker; {
				case lengthLeftInSubArrForMoreString >= 1:
					// return the current index
					return i
					println("we are in the section where there is lengthLeftInSubArrForMoreString >= 1 is true, btw the full subtitle in the row is ->", transcript.Subtitles[i].Text)
				case lengthLeftInSubArrForMoreString < 1:
					println("we are in the section  lengthLeftInSubArrForMoreString<1 and the full subtitles is ->", transcript.Subtitles[i].Text)
					// there is no base case by(mathematically) btw
					if len(transcript.Subtitles) > i+1 {
						return i + 1
					} else {
						return i
					}
				default:
					println("----\n\n\n\n  we hit the base case in the switch statement this should be not happening(mathematically) fix this    ----\n\n\n\n")
					return i
				}
			}
			if sponserShipIndexTracker > sponserShipAt {
				println("the sponserShipAt(index):->", sponserShipAt, " the current index of words is (sponserShipIndexTracker, track of if we reached sponserShipAt)  ", sponserShipIndexTracker, " and the index found at(rough one) the sponserShipIndexInSubtitlesArray.IndexAt", sponserShipIndexInSubtitlesArray.IndexAt, "word at:", j, " at is ->", wordWithoutSpace, "--")
				println("--this one is less desirable ---")
				println("but fuck it had to return the current Subtitle")
				return i
			}
			println("word at:", j, " at is ->", wordWithoutSpace, "--", "sponserShipIndexTracker:", sponserShipIndexTracker, " and sponserShi starts at ->", sponserShipAt)
			sponserShipIndexTracker += len(" ")
		}
		println("+++++++++++++++++++++++++++++++++++++++++++++++++ \n")
		sponserShipIndexTracker += len(" ")
	}
	println("in the end of getPreciseIndexOfSubtitle() and the sponserShipAt:->", sponserShipAt, " and the index found at(rough one) the sponserShipIndexInSubtitlesArray.IndexAt", sponserShipIndexInSubtitlesArray.IndexAt)
	println(" -- you should probally not reach here but since you have I will just accept that there is some error and will return the index  ")
	return sponserShipIndexInSubtitlesArray.IndexAt
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
