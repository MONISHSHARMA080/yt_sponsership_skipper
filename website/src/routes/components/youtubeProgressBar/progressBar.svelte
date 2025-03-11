
<script lang="ts">
     let {sponsorStart, sponsorEnd, videoLength= sponsorStart + 30, baseTailwindBG = "bg-gray-500", progressTailwindBG = "bg-red-500", sponsorTailwindBG = "bg-yellow-400",
      funcToRunWhenOnTheSponSorSection, playVideo = true, funcToRunAfterTheSponsorSegment
     }:
    //  props type
     { sponsorStart:number, sponsorEnd:number, videoLength?:number, baseTailwindBG?:string, progressTailwindBG?:string, sponsorTailwindBG?:string
      ,funcToRunWhenOnTheSponSorSection:(areWeInSponsorSegment:Boolean)=>void, playVideo:boolean, funcToRunAfterTheSponsorSegment?:()=>void
      } = $props()

     let progress = $state(0);
     let progressPercentage = $derived((progress / videoLength) * 100);
     let isPlaying = $derived(playVideo);

  let progressInterval:any;
  
  $effect(() => {
    if (isPlaying) {
		
      progressInterval = setInterval(() => {
        progress += 0.1;
        
        if (progress >= videoLength) {
          // isPlaying = false;
            playVideo = false;
          progress = 0;
        }
      }, 10);
    } else if (progressInterval) {
      clearInterval(progressInterval);
    }

    return () => {
      if (progressInterval) clearInterval(progressInterval);
    };
  });


 let inSponsorSegment = $derived(progress >= sponsorStart && progress <= sponsorEnd)

  $effect(() => {
    if(inSponsorSegment){
      funcToRunWhenOnTheSponSorSection(inSponsorSegment)
    }
  })
  
 let afterTheSponsorSegment = $derived(progress > sponsorEnd)
 $effect(()=>{
  if (afterTheSponsorSegment && funcToRunAfterTheSponsorSegment ) {
    funcToRunAfterTheSponsorSegment()
  }
 })
  function togglePlay() {
	 playVideo= !playVideo;
  }

  
</script>


<div class="absolute bottom-0 left-0 right-0 p-4 ">
        <!-- Progress bar -->
        <div class="h-1 w-full {baseTailwindBG} rounded-full mb-4 overflow-hidden relative">
          <!-- Main progress - Using progressPercentage instead of progress -->
          <div class="h-full {progressTailwindBG} rounded-4xl" style="width: {progressPercentage}%"></div>
          
          <!-- Sponsorship marker -->
          <div class="absolute top-0 h-1 z-30 {sponsorTailwindBG}  rounded-3xl"
             style="left: {sponsorStart}%; width: {sponsorEnd - sponsorStart}%;" 
          ></div>
        </div>
</div>