var server = "localhost";

var leadBook = "";
var leadChapter = "";
var leadVerse = "";

document.getElementById("searchBox").onkeydown = function(event) {
	if (event.keyCode == 13) {
		loadVerses()
	};
};


function loadVerses() {
  var queryText = document.getElementById("searchBox").value
  var verseDiv = document.createElement("div");
  verseDiv.className = "verse_box";
  var verseNumDiv = document.createElement("div");
  verseNumDiv.className = "verse_number";
  var verseTextDiv = document.createElement("div");
  verseTextDiv.className = "verse_text";

  var verseNumP = document.createElement("p");
  var verseTextP = document.createElement("p");

  var ajaxRequest = new XMLHttpRequest();
  var apiRequest = "http://" + server + "/v0/query?query=" + queryText
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	var responseData = JSON.parse(ajaxRequest.responseText);
	leadBook = responseData['Reference']['Book']
	leadChapter = responseData['Reference']['Chapter']
	leadVerse = responseData['Reference']['VerseNumber']

	verseNumP.innerHTML = responseData['Reference']['Chapter'] + ":" + responseData['Reference']['VerseNumber'];
	verseTextP.innerHTML = responseData['Text'];

	verseNumDiv.appendChild(verseNumP);
	verseTextDiv.appendChild(verseTextP);
  	verseDiv.appendChild(verseNumDiv);
  	verseDiv.appendChild(verseTextDiv);
	var verseDivs = document.getElementsByClassName("verse_box");
	if (verseDivs) {
	  while(verseDivs[0]) {
		  verseDivs[0].parentNode.removeChild(verseDivs[0]);
	  }
	}
	document.getElementById("content").appendChild(verseDiv);
	
	document.getElementById("simpleInfo").innerHTML = responseData['Reference']['Book']
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};

function previousVerse() {
  var verseDiv = document.createElement("div");
  verseDiv.className = "verse_box";
  var verseNumDiv = document.createElement("div");
  verseNumDiv.className = "verse_number";
  var verseTextDiv = document.createElement("div");
  verseTextDiv.className = "verse_text";

  var verseNumP = document.createElement("p");
  var verseTextP = document.createElement("p");

  var ajaxRequest = new XMLHttpRequest();
  var apiRequest = "http://" + server + "/v0/previous?ref=" + leadBook + "~" + leadChapter + "." + leadVerse;
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	var responseData = JSON.parse(ajaxRequest.responseText);
	leadBook = responseData['Reference']['Book']
	leadChapter = responseData['Reference']['Chapter']
	leadVerse = responseData['Reference']['VerseNumber']

	verseNumP.innerHTML = responseData['Reference']['Chapter'] + ":" + responseData['Reference']['VerseNumber'];
	verseTextP.innerHTML = responseData['Text'];

	verseNumDiv.appendChild(verseNumP);
	verseTextDiv.appendChild(verseTextP);
  	verseDiv.appendChild(verseNumDiv);
  	verseDiv.appendChild(verseTextDiv);
	var verseDivs = document.getElementsByClassName("verse_box");
	if (verseDivs) {
	  while(verseDivs[0]) {
		  verseDivs[0].parentNode.removeChild(verseDivs[0]);
	  }
	}
	document.getElementById("content").appendChild(verseDiv);
	
	document.getElementById("simpleInfo").innerHTML = responseData['Reference']['Book']
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
	
};
