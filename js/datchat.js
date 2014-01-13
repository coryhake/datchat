function storePhone() {
  var phone = document.getElementById("phone");
  alert("Phone stored as: " + phone.value);
}
    
function longpoll(url, callback) {

	var req = new XMLHttpRequest (); 
	req.open ('GET', url, true); 

	req.onreadystatechange = function (aEvt) {
		if (req.readyState == 4) { 
			if (req.status == 200) {
				callback(req.responseText);
				longpoll(url, callback);
			} else {
				alert ("long-poll connection lost");
			}
		}
	};

	req.send(null);
}

function recv(msg) {

	var box = document.getElementById("box");
	
	box.value += "\n" + msg;
}

function send() {

	var from = document.getElementById("from");
	var rcpt = document.getElementById("rcpt");
	var box = document.getElementById("box");
	var input = document.getElementById("input");

	var req = new XMLHttpRequest (); 
  req.open ('POST', "/push?rcpt=" + rcpt.value, true);
	req.onreadystatechange = function (aEvt) {
		if (req.readyState == 4) { 
			if (req.status == 200) {
			} else {
				alert ("failed to send!");
			}
		}
	};

	req.send(from.value + ": " + input.value);
	
	box.value += "\nme: " + input.value;
	input.value = "";
}

function login() {
 
	var from = document.getElementById("from");
  if(from.value != "") {
    longpoll("/poll?rcpt=" + from.value, recv);
    alert("Logged in as: " + from.value);
  }
  else{
   alert("Display name cannot be blank!"); 
  }
}