package askllm

import (
	"strconv"
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

func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming, logger *zap.Logger) {
	// func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponserSubtitleFromLLM *string, full_captions *string, logger *zap.Logger)
	subtitleTimming := ResponseForGettingSubtitlesTiming{}
	newFullCaptions := SanatizeStrignsForSearching(full_captions)
	newsponserSubtitleFromLLM := SanatizeStrignsForSearching(sponserSubtitleFromLLM)
	full_captions = &newFullCaptions
	sponserSubtitleFromLLM = &newsponserSubtitleFromLLM
	logger.Info("got the full caption in the GetTimeAndDurInTheSubtitles2 and", zap.String("*full_captions is", *full_captions))
	logger.Info("got the sponsership_subtitles_form_groq in the GetTimeAndDurInTheSubtitles2 and", zap.String("*sponserSubtitleFromLLM is", *sponserSubtitleFromLLM))

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
	logger.Info("some var info at the start of the GetTimeAndDurInTheSubtitles2 ->", zap.Int("sponsershipStartIndex is ->", sponsershipStartIndex),
		zap.Int("sponserShipEndIndex is ->", sponserShipEndIndex), zap.Bool("is the sponserSubtitleFromLLM == full_captions[sponsershipStartIndex : sponserShipEndIndex ]", *sponserSubtitleFromLLM == (*full_captions)[sponsershipStartIndex:sponserShipEndIndex]))

	sponserShipStartIndexForSubtitlesArray := sponsershipPositionIndex{Found: false, IndexAt: 0}
	sponserShipEndIndexForSubtitlesArray := sponsershipPositionIndex{Found: false, IndexAt: 0}

	strInSubtitlesLenCounter := 0
	prevSubLenCounter := 0
	secondLastLenCounter := 0
	for i, subtitle := range transcripts.Subtitles {

		secondLastLenCounter = prevSubLenCounter
		prevSubLenCounter = strInSubtitlesLenCounter

		strInSubtitlesLenCounter += len(subtitle.Text + " ")
		// if we are after the sponsor ship start index and it is not assigned before
		if !sponserShipStartIndexForSubtitlesArray.Found && strInSubtitlesLenCounter >= sponsershipStartIndex {
			sponserShipStartIndexForSubtitlesArray.Found = true
			sponserShipStartIndexForSubtitlesArray.IndexAt = i
			sponserShipStartIndexForSubtitlesArray.SubtitleStringCounterAtTheCurrentIteration = strInSubtitlesLenCounter
			sponserShipStartIndexForSubtitlesArray.SubtitleStringCounterAtThePreviosLoopIteration = prevSubLenCounter
			sponserShipStartIndexForSubtitlesArray.SubtitleStringCounterAtThe2ndPreviousLoopIteration = secondLastLenCounter
			logger.Info("we have reached the sponserShipStartIndex ", zap.Int("the index is ", i), zap.String("text at that index in the subtitle is ->", subtitle.Text))
		}
		if !sponserShipEndIndexForSubtitlesArray.Found && strInSubtitlesLenCounter >= sponserShipEndIndex {
			sponserShipEndIndexForSubtitlesArray.Found = true
			sponserShipEndIndexForSubtitlesArray.IndexAt = i
			sponserShipEndIndexForSubtitlesArray.SubtitleStringCounterAtTheCurrentIteration = strInSubtitlesLenCounter
			sponserShipEndIndexForSubtitlesArray.SubtitleStringCounterAtThePreviosLoopIteration = prevSubLenCounter
			sponserShipEndIndexForSubtitlesArray.SubtitleStringCounterAtThe2ndPreviousLoopIteration = secondLastLenCounter
			logger.Info("we have reached the sponserShipEndIndex ", zap.Int("the index is ", i), zap.String("text at that index in the subtitle is ->", subtitle.Text))
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
	preciseStartIndex := getPreciseIndexOfSubtitle(transcripts, &sponserShipStartIndexForSubtitlesArray, sponsershipStartIndex, logger)
	logger.Info("============= not getting the precise end index subtitle================")
	preciseEndIndex := getPreciseIndexOfSubtitle(transcripts, &sponserShipEndIndexForSubtitlesArray, sponserShipEndIndex, logger)

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

	logger.Info("got the both precise start time and the end time , logging both of them ", zap.Int("precise start index is", preciseStartIndex), zap.String("text at the precise start Index is ", transcripts.Subtitles[preciseStartIndex].Text),
		zap.Int("precise end index is", preciseEndIndex), zap.String("text at the preciseEndIndex is ", transcripts.Subtitles[preciseEndIndex].Text))
	startTime, err := getTimeOutOfString(transcripts.Subtitles[preciseStartIndex].Start)
	if err != nil {
		logger.Error("the error in the startTime's precise index str to time(get startTime  from the transcripts and parse it into float), now we will try with the less accurate time we got from the first loop(not look at the adjacent index and edge cases) ",
			zap.Error(err))
		startTime, err := getTimeOutOfString(transcripts.Subtitles[sponserShipStartIndexForSubtitlesArray.IndexAt].Start)
		if err != nil {
			logger.Error("the error in the startTime's sponserShipStartIndexForSubtitlesArray.IndexAt index str to time conversion, we cannot get the time even from second time after first error, sendign error response",
				zap.Error(err))
			subtitleTimming.Err = err
			responseForTimmingChannel <- subtitleTimming
			return
		}
		subtitleTimming.StartTime = int(startTime)
		logger.Info("got the precise index for the start time ", zap.Int("preciseStartIndex is ", preciseStartIndex))
	} else {
		subtitleTimming.StartTime = int(startTime)
		logger.Info("got the precise index for the start time ", zap.Int("preciseStartIndex is ", preciseStartIndex))
	}

	endTime, err := getTimeOutOfString(transcripts.Subtitles[preciseEndIndex].Start)
	if err != nil {
		logger.Error("the error in the endTime's precise index str to time(get endTime  from the transcripts and parse it into float), now we will try with the less accurate time we got from the first loop(not look at the adjacent index and edge cases) ", zap.Error(err))
		endTime, err := getTimeOutOfString(transcripts.Subtitles[sponserShipEndIndexForSubtitlesArray.IndexAt].Start)
		if err != nil {
			logger.Error("the error in the endTime's sponserShipStartIndexForSubtitlesArray.IndexAt index str to time conversion, we cannot get the time even from second time after first error, sendign error response",
				zap.Error(err))
			subtitleTimming.Err = err
			responseForTimmingChannel <- subtitleTimming
			return
		}
		subtitleTimming.EndTime = int(endTime)
		logger.Info("got the precise index for the end time ", zap.Int("preciseEndIndex is ", preciseEndIndex))
	} else {
		subtitleTimming.EndTime = int(endTime)
		logger.Info("got the precise index for the end time ", zap.Int("preciseEndIndex is ", preciseEndIndex))
	}
	subtitleTimming.Err = nil
	responseForTimmingChannel <- subtitleTimming
}

// this func is designed to be used when we have gotten the subtitle but it there is a possibility that it could be a error as it might be in the adjacent index
func getPreciseIndexOfSubtitle(transcript *Transcripts, sponserShipIndexInSubtitlesArray *sponsershipPositionIndex, sponserShipAt int, logger *zap.Logger) int {
	// this is a  generic, so it can be used as the first one and the last one
	// ok get the index and form the index-1 to index +1 if we are not able to get it back then we return error
	//
	// the sponserShipAt is the position where the subtitle starts or ends
	// +2 here as we want the range to be i-1 , i , i+1
	logger.Info("in the getPreciseIndexOfSubtitle() and here is the assertion log", zap.Bool("we assert that the sponserShipAt< sponserShipIndexInSubtitlesArray.IndexAt", sponserShipAt < sponserShipIndexInSubtitlesArray.IndexAt))

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
				logger.Info(" in the case where sponserShipIndexTracker == sponserShipAt and we got the exact position")
				logger.Info("-the sponserShipAt(index):->",
					zap.Int("sponserShipAt", sponserShipAt),
					zap.String("message", " the current index of words is (sponserShipIndexTracker, track of if we reached sponserShipAt)  "),
					zap.Int("sponserShipIndexTracker", sponserShipIndexTracker),
					zap.String("message2", " and the index found at(rough one) the sponserShipIndexInSubtitlesArray.IndexAt"),
					zap.Int("IndexAt", sponserShipIndexInSubtitlesArray.IndexAt),
					zap.String("word_info", "word at:"),
					zap.Int("j", j),
					zap.String("word", " at is ->"+wordWithoutSpace+"--   "))

				// the reason for this if statement is that the  sponserShipIndexTracker == sponserShipAt means we are only one index(in strings) behind the sponserShip start and

				// in the edge case (assume last word in the end of the subtitles index) the sponseShip is in the next index and to be presice we need to skip over
				// but in the case it is for the end function there the end index is the end. These differences is cause of how the len fuc work and how we decided to calculate the
				// start index and the end one(strings.Index gived you the start but the index + full_captions gives you the end word)
				//
				// ya we ditched that cause if the end sponserShip is at the subtitle end then we have a problem as the video should start from the next index and
				// so we are doing that
				switch lengthLeftInSubArrForMoreString := len(transcript.Subtitles[i].Text) - sponserShipIndexTracker; {
				case lengthLeftInSubArrForMoreString >= 1:
					// return the current index
					logger.Info("returning current subtitle index as precise match and the length left in the sub array after - sponserShipIndexTracker is >=1", zap.Int("index", i))
					return i
				case lengthLeftInSubArrForMoreString < 1:
					// there is no base case by(mathematically) btw
					if len(transcript.Subtitles) > i+1 {
						logger.Info("overflow to next subtitle index", zap.Int("nextIndex", i+1))
						return i + 1
					} else {
						logger.Info("no next subtitle, returning current index", zap.Int("index", i))
						return i
					}
				default:
					logger.Warn("unexpected switch default case reached this should not have been mathematically possibile", zap.Int("index", i))
					return i
				}
			}
			if sponserShipIndexTracker > sponserShipAt {
				logger.Info("tracker surpassed target, returning current index --this one is less desirable --- but fuck it had to return the current Subtitle ",
					zap.Int("targetIndex or the the sponserShipAt(index):->", sponserShipAt),
					zap.Int("trackerIndex or the current index of words is (sponserShipIndexTracker, track of if we reached sponserShipAt) ", sponserShipIndexTracker),
					zap.Int("index found at(rough one/ backup / the one we we going to replacing with the precise one ) the sponserShipIndexInSubtitlesArray.IndexAt ", sponserShipIndexInSubtitlesArray.IndexAt),
				)
				return i
			}
			logger.Info("reached the end of the loop over the words in the array wihout whitespace, at a index, here is some info ", zap.Int("at index", j),
				zap.String("word here is ->", wordWithoutSpace), zap.Int("sponserShipIndexTracker is ", sponserShipIndexTracker),
				zap.Int("sponserShi starts at", sponserShipAt),
			)
			sponserShipIndexTracker += len(" ")
		}
		logger.Info("+++++++++++++++++++++++++++++++++++++++++++++++++ ")
		sponserShipIndexTracker += len(" ")
	}
	logger.Info("in the end of getPreciseIndexOfSubtitle() and the sponserShipAt:->",
		zap.Int("sponserShipAt", sponserShipAt),
		zap.String("message", " and the index found at(rough one) the sponserShipIndexInSubtitlesArray.IndexAt"),
		zap.Int("IndexAt", sponserShipIndexInSubtitlesArray.IndexAt))
	logger.Info(" -- you should probally not reach here but since you have I will just accept that there is some error and will return the index  ")
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

func getTimeOutOfString(t string) (float64, error) {
	i, err := strconv.ParseFloat(t, 64)
	return i, err
}
