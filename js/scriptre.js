var server = "localhost";

var leadBook = "";
var leadChapter = "";
var leadVerse = "";
var responseData;

document.getElementById("searchBox").onkeydown = function(event) {
	if (event.keyCode == 13) {
		loadVerses()
	};
};

function loadVerses() {
  
  var queryText = document.getElementById("searchBox").value;
  getVerses(queryText);

};

function clearContent() {
  var verseDivs = document.getElementsByClassName("verse_box");
  if (verseDivs) {
    while(verseDivs[0]) {
      verseDivs[0].parentNode.removeChild(verseDivs[0]);
    };
  };
};

function getVerses(query) {
  var ajaxRequest = new XMLHttpRequest();
  var apiRequest = "http://" + server + "/v0/query?query=" + query
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	responseData = JSON.parse(ajaxRequest.responseText);
	clearContent()
	makeVerseBox(responseData);
  	document.getElementById("simpleInfo").innerHTML = responseData['Reference']['Book'];
	leadBook = responseData['Reference']['Book']
	leadChapter = responseData['Reference']['Chapter']
	leadVerse = responseData['Reference']['VerseNumber']
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};

function makeVerseBox(verseData) {
  var verseDiv = document.createElement("div");
  verseDiv.className = "verse_box";
  var verseNumDiv = document.createElement("div");
  verseNumDiv.className = "verse_number";
  var verseTextDiv = document.createElement("div");
  verseTextDiv.className = "verse_text";

  var verseNumP = document.createElement("p");
  var verseTextP = document.createElement("p");

  verseNumP.innerHTML = verseData['Reference']['Chapter'] + ":" + verseData['Reference']['VerseNumber'];
  verseTextP.innerHTML = verseData['Text'];

  verseNumDiv.appendChild(verseNumP);
  verseTextDiv.appendChild(verseTextP);
  verseDiv.appendChild(verseNumDiv);
  verseDiv.appendChild(verseTextDiv);
  document.getElementById("content").appendChild(verseDiv);
};


function previous(mode) {
  var ajaxRequest = new XMLHttpRequest();
  if (!mode) {
	mode = "verse"
	var apiRequest = "http://" + server + "/v0/previous?ref=" + leadBook + "~" + leadChapter + "." + leadVerse;
  } else if (mode == "chapter") {
	var apiRequest = "http://" + server + "/v0/previous?ref=" + leadBook + "~" + leadChapter + "." + leadVerse + "&mode=chapter";
  };
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	responseData = JSON.parse(ajaxRequest.responseText);
	clearContent()
	makeVerseBox(responseData);
  	document.getElementById("simpleInfo").innerHTML = responseData['Reference']['Book'];
	leadBook = responseData['Reference']['Book']
	leadChapter = responseData['Reference']['Chapter']
	leadVerse = responseData['Reference']['VerseNumber']
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};

function next(mode) {
  var ajaxRequest = new XMLHttpRequest();
  if (!mode) {
	mode = "verse"
	var apiRequest = "http://" + server + "/v0/next?ref=" + leadBook + "~" + leadChapter + "." + leadVerse;
  } else if (mode == "chapter") {
	var apiRequest = "http://" + server + "/v0/next?ref=" + leadBook + "~" + leadChapter + "." + leadVerse + "&mode=chapter";
  };
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	responseData = JSON.parse(ajaxRequest.responseText);
	clearContent()
	makeVerseBox(responseData);
  	document.getElementById("simpleInfo").innerHTML = responseData['Reference']['Book'];
	leadBook = responseData['Reference']['Book']
	leadChapter = responseData['Reference']['Chapter']
	leadVerse = responseData['Reference']['VerseNumber']
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};

