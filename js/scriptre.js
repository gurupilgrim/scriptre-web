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
	responseData.forEach(function(item,key){
		makeVerseBox(item);
	});
  	document.getElementById("simpleInfo").innerHTML = responseData[Object.keys(responseData)[0]]['Reference']['Book'];
	leadBook = responseData[Object.keys(responseData)[0]]['Reference']['Book']
	leadChapter = responseData[Object.keys(responseData)[0]]['Reference']['Chapter']
	leadVerse = responseData[Object.keys(responseData)[0]]['Reference']['VerseNumber']
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
  verseDiv.id = verseData['Reference']['Chapter'] + "-" + verseData['Reference']['VerseNumber'];
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

function prependVerseBox(verseData) {
  var verseDiv = document.createElement("div");
  verseDiv.className = "verse_box";
  verseDiv.id = verseData['Reference']['Chapter'] + "-:" + verseData['Reference']['VerseNumber'];
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
  contentDiv = document.getElementById("content");
  contentDiv.insertBefore(verseDiv, contentDiv.childNodes[0]);
}

function previous(mode) {
  var ajaxRequest = new XMLHttpRequest();
  if (!mode) {
	if (leadVerse < 2) {
		return
	}
	mode = "verse"
	var apiRequest = "http://" + server + "/v0/previous?ref=" + leadBook + "~" + leadChapter + "." + leadVerse;
  } else if (mode == "chapter") {
	if (leadChapter < 2) {
		return
	}
	var apiRequest = "http://" + server + "/v0/previous?ref=" + leadBook + "~" + leadChapter + "." + leadVerse + "&mode=chapter";
  };
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	responseData = JSON.parse(ajaxRequest.responseText);
	if (mode == "verse") {
		prependVerseBox(responseData);
  		document.getElementById("simpleInfo").innerHTML = responseData['Reference']['Book'];
		leadBook = responseData['Reference']['Book']
		leadChapter = responseData['Reference']['Chapter']
		leadVerse = responseData['Reference']['VerseNumber']
	}
	if (mode == "chapter") {
		clearContent()
		responseData.forEach(function(item,key){
			makeVerseBox(item);
		});
		leadBook = responseData[Object.keys(responseData)[0]]['Reference']['Book']
		leadChapter = responseData[Object.keys(responseData)[0]]['Reference']['Chapter']
		leadVerse = responseData[Object.keys(responseData)[0]]['Reference']['VerseNumber']
	}
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};

function next(mode) {
  if (!mode) {
	mode = "verse"
  } else if (mode == "chapter") {
	var apiRequest = "http://" + server + "/v0/next?ref=" + leadBook + "~" + leadChapter + "." + leadVerse + "&mode=chapter";
  };
  if (mode == "verse") {
	var contentDiv = document.getElementById("content");
	verseBoxInfo = getCurrentVerseBox();
	verseBoxInfoStrings = verseBoxInfo['id'].split("-");
	verseBoxInfoStrings[verseBoxInfoStrings.length - 1] = parseInt(verseBoxInfoStrings[verseBoxInfoStrings.length -1]) + 1;
	verseBoxInfo['id'] = verseBoxInfoStrings.join("-");
	nextVerseBox = document.getElementById(verseBoxInfo['id']);
	console.log(nextVerseBox);
	console.log(nextVerseBox.offsetTop);
	// scroll to that id
	contentDiv.scrollTop = nextVerseBox.offsetTop - contentDiv.offsetTop;
	// set the lead reference
	leadVerse = leadVerse + 1;
	return
  }
  var ajaxRequest = new XMLHttpRequest();
  ajaxRequest.open("GET", apiRequest, true);

  ajaxRequest.onload = function() {
	responseData = JSON.parse(ajaxRequest.responseText);
	clearContent()
	responseData.forEach(function(item,key){
		makeVerseBox(item);
	});
  	document.getElementById("simpleInfo").innerHTML = responseData[Object.keys(responseData)[0]]['Reference']['Book'];
	leadBook = responseData[Object.keys(responseData)[0]]['Reference']['Book']
	leadChapter = responseData[Object.keys(responseData)[0]]['Reference']['Chapter']
	leadVerse = responseData[Object.keys(responseData)[0]]['Reference']['VerseNumber']
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  };
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};

function getCurrentVerseBox() {
	var contentDiv = document.getElementById("content")
	var contentChildren = Array.prototype.slice.call(contentDiv.childNodes);
	var result = new Array();
	contentChildren.some(function(child){
	//[].forEach.call(contentChildren, function(child) {
		if ((contentDiv.scrollTop + contentDiv.offsetTop) <= child.offsetTop) {
			result['id'] = child.id;
			result['offset'] = child.offsetTop;
			return result;
		}
	});
	return result;
}
