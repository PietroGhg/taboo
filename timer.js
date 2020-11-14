//calculate the seconds
var secs = 60;
var count = secs;

document.getElementById("seconds").value = secs;

function countdown() {
    running_time = true;
    Decrement();
}

function Decrement(){
    console.log(count);
    if(count <= 1){
	running_time = false;
	count = secs;
	document.getElementById("seconds").value = secs;
	alert("Tempo scaduto");
    }
    else{
	count--;
	document.getElementById("seconds").value = count;
	setTimeout('Decrement()', 1000);
    }
}
	
	



