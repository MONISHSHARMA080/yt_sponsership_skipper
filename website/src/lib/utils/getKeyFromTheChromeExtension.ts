export function getTheKeyFromTheEvent(){
    console.log("in the event function");
    
 document.addEventListener('getKeyForTheSvelteWebside',(event)=>{
    console.log("the event received and it is ->",event);
 })
}