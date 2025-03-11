<script lang="ts">
	import { keyFromChromeExtensionState } from '$lib/sharedState/sharedKeyState.svelte';
	import { interactWithTheChromeExtensionAndStoreIt } from '$lib/utils/interactWithChromeExtension/getKeyFromTheChromeExtension';
	import { sendChromeExtensionNewKey } from '$lib/utils/interactWithChromeExtension/sendNewKeyAfterPayment';
	import { onMount } from 'svelte';
	import Component14 from './components/component14.svelte';
	import { Tween } from 'svelte/motion';
	import { cubicIn, cubicInOut, cubicOut } from 'svelte/easing';
	import { Play, FastForward, Pause } from 'lucide-svelte';
	import { fade } from 'svelte/transition';
	import ProgressBar from './components/youtubeProgressBar/progressBar.svelte';
	
	let change = $state(false);
	
	onMount(() => {
		setTimeout(() => {
			change = true;
		console.log("the event is running ->");
		}, 100);
		
		// Commented extension code preserved as in original
		// let interactWithExtensionClass = new interactWithTheChromeExtensionAndStoreIt
		// let error = interactWithExtensionClass.start((key)=>{console.log("the key is received and it is ->",key," --- about to update the svelete store")
		// 	// keyFromChromeExtensionState.key = key
		// 	interactWithExtensionClass.cleanup()
		// // let checkKeyAndnew = new checkIfKeyIsValidAndUpdateTheState()
		// //  checkKeyAndnew.seeIfKeyIsValid(key)
		// }
		// console.log("error in interacting with the chrome extension is -> ",error );
		// const sendNewKeyClass = new sendChromeExtensionNewKey("(((((((((((((((((((((((((((((((")
		// sendNewKeyClass.sendKey().then((response)=>{
		// 	console.log("the response after sending the key is  ->",response)
		// 	console.log('\n\n\n\n');
		// 	// sendNewKeyClass.clearEventListener()
		// })
	});
	

  let isPlaying = $state(false);
  let progress = $state(0);
  
  // Define the sponsor segment (in real app, this would come from an API or database)
  let sponsorStart = $state(10);
  let sponsorEnd = $state(40);
  let videoLength = $state(50);
  
  // Interval for updating progress
  let progressInterval:any;
  
  $effect(() => {
    if (isPlaying) {
		
      progressInterval = setInterval(() => {
        progress += 0.1;
        
        if (progress >= videoLength) {
          isPlaying = false;
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
  
  // Determine if we're currently in a sponsor segment
  let inSponsorSegment = $derived(progress >= sponsorStart && progress <= sponsorEnd)
  
  function togglePlay() {
	isPlaying = !isPlaying;
  console.log(`is playing -> ${isPlaying}`);
  
  }

  // Calculate progress as percentage of total video length
  let progressPercentage = $derived((progress / videoLength) * 100);
  
</script>

{#if change}
 
<button onclick={togglePlay} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
play the video
</button>
 

<ProgressBar
  funcToRunWhenOnTheSponSorSection={(areWeInSponsorSegment:Boolean)=>{console.log("are we in the sponsor segment ->",areWeInSponsorSegment)}}
  sponsorStart={sponsorStart}
  sponsorEnd={sponsorEnd}
  videoLength={videoLength}
  playVideo={isPlaying}
  funcToRunAfterTheSponsorSegment={()=>{console.log(" sponsor segment ended and stopping it ")
  isPlaying = false
}}


/>

      <!-- Video controls -->
      <!-- <div class="absolute bottom-0 left-0 right-0 p-4 "> -->
        <!-- Progress bar -->
        <!-- <div class="h-1 w-full bg-gray-700 rounded-full mb-4 overflow-hidden relative"> -->
          <!-- Main progress - Using progressPercentage instead of progress -->
          <!-- <div class="h-full bg-red-500 rounded-4xl" style="width: {progressPercentage}%"></div> -->
          
          <!-- Sponsorship marker -->
          <!-- <div class="absolute top-0 h-1 z-30 bg-yellow-400 rounded-3xl"
             style="left: {sponsorStart}%; width: {sponsorEnd - sponsorStart}%;" 
          ></div> -->
        <!-- </div> -->
	  <!-- </div> -->

      



   
{:else}
    <h1>Welcome to SvelteKit</h1>
    {#if keyFromChromeExtensionState.key === null}
        <h2>bro the key is not there ->{keyFromChromeExtensionState.key}</h2>
    {:else}
        <h2>the key is -> {keyFromChromeExtensionState.key}</h2>
    {/if}
    {#if keyFromChromeExtensionState.isValidatedThroughBackend}
        <h1>Hello {keyFromChromeExtensionState.name} and your email is {keyFromChromeExtensionState.email}</h1>
    {/if}
    <p >
        the key value it {JSON.stringify(keyFromChromeExtensionState)}
    </p>
{/if}