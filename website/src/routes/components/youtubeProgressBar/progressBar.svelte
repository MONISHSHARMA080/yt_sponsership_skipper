<script lang="ts">
	import { PowerOff } from "lucide-svelte";
     let {
        sponsorStart, 
        sponsorEnd, 
        videoLength = 100, // Default value should be independent of sponsorStart
        baseTailwindBG = "bg-gray-500", 
        progressTailwindBG = "bg-red-500", 
        sponsorTailwindBG = "bg-yellow-400",
        funcToRunWhenInTheSponSorSection,
        playVideo = true, 
        funcToRunAfterTheSponsorSegment,
        loopEnabled = false, // Added loop functionality
        autoRestartAfterCompletion = true, // Added option to auto restart
        funcToRunAfterVideoCompletion
     }: { 
        sponsorStart: number, 
        sponsorEnd: number, 
        videoLength?: number, 
        baseTailwindBG?: string, 
        progressTailwindBG?: string, 
        sponsorTailwindBG?: string,
        funcToRunWhenInTheSponSorSection: (areWeInSponsorSegment: boolean) => void, 
        playVideo: boolean, 
        funcToRunAfterTheSponsorSegment?: () => void,
        loopEnabled?: boolean, // Added type definition
        autoRestartAfterCompletion?: boolean // Added type definition
        funcToRunAfterVideoCompletion?: () => void
     } = $props();
     
     let progress = $state(0);
     let progressPercentage = $derived((progress / videoLength) * 100);
     let isPlaying = $derived(playVideo);
     let progressInterval: any;
     
     // Calculate sponsor percentages for display
     let sponsorStartPercent = $derived((sponsorStart / videoLength) * 100);
     let sponsorEndPercent = $derived((sponsorEnd / videoLength) * 100);
     // so that we only call it once
     let calledTheFuncToRunAfterSponsorBeofre = $state(false);
     let completionCount = $state(0); // Track number of complete playbacks
  
     function resetProgress() {
       progress = 0;
       calledTheFuncToRunAfterSponsorBeofre = false;
     }
  
     function togglePlay() {
       playVideo = !playVideo;
     }
  
     function restartVideo() {
       resetProgress();
       playVideo = true;
     }
  
     $effect(() => {
       if (isPlaying) {
         progressInterval = setInterval(() => {
           progress += 0.099;
           
           if (progress >= videoLength) {
             completionCount++;
             
             if (loopEnabled || autoRestartAfterCompletion) {
               // If loop is enabled, just reset progress
               resetProgress();
               // If autoRestart is not enabled with loop, pause playback
               if (!autoRestartAfterCompletion) {
                 playVideo = false;
               }
             } else {
               // Standard behavior - stop at the end
               playVideo = false;
               progress = 0;
             }
           }
         }, 10);
       } else if (progressInterval) {
         clearInterval(progressInterval);
       }
       return () => {
         if (progressInterval) clearInterval(progressInterval);
       };
     });
  
     // - 1 so that we can start before and get the animation right
     let inSponsorSegment = $derived(progress >= sponsorStart -1.6 && progress <= sponsorEnd);
     
     $effect(() => {
       if(inSponsorSegment) {
         funcToRunWhenInTheSponSorSection(inSponsorSegment);
       }
     });
  
     let afterTheSponsorSegment = $derived(progress > sponsorEnd && calledTheFuncToRunAfterSponsorBeofre === false);
     
     $effect(() => {
       if (afterTheSponsorSegment && funcToRunAfterTheSponsorSegment) {
         // Call the callback function if provided
         funcToRunAfterTheSponsorSegment();
         calledTheFuncToRunAfterSponsorBeofre = true;
         
         console.log(`Progress: ${progress}, Video Length: ${videoLength}, Sponsor End: ${sponsorEnd}, Sponsor Start: ${sponsorStart}`);
       }
     });

     let afterVideoCompletion = $derived(progress >= videoLength -2 );
     $effect(()=>{
      
        if(afterVideoCompletion && funcToRunAfterVideoCompletion){
          funcToRunAfterVideoCompletion();
        }
     })


</script>

<div class="absolute bottom-0 left-0 h-4 rounded my w-full">
  <!-- Progress bar -->
  <div class="h-4 w-full {baseTailwindBG} rounded mb-4 overflow-hidden relative">
    <!-- Main progress -->
    <div class="h-4 {progressTailwindBG} rounded" style="width: {progressPercentage}%"></div>
    
    <!-- Sponsorship marker - Using calculated percentages -->
    <div class="absolute top-0 h-4 z-30 {sponsorTailwindBG} rounded"
         style="left: {sponsorStartPercent}%; width: {sponsorEndPercent - sponsorStartPercent}%;" 
    ></div>
  </div>
</div>