var server = "script.re"

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
	verseNumP.innerHTML = responseData['Reference']['Book'] + " " + responseData['Reference']['Chapter'] + ":" + responseData['Reference']['VerseNumber'];
	verseTextP.innerHTML = responseData['Text'];

	verseNumDiv.appendChild(verseNumP);
	verseTextDiv.appendChild(verseTextP);
  	verseDiv.appendChild(verseNumDiv);
  	verseDiv.appendChild(verseTextDiv);
	document.getElementById("content").appendChild(verseDiv);
  };
  ajaxRequest.onerror = function() {
  	alert("Error!");
  }
  //ajaxRequest.responseType = 'json';
  ajaxRequest.send();
  console.log(ajaxRequest);
};
