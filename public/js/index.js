// region Relative Time Formatting garbage
// https://stackoverflow.com/questions/6108819/javascript-timestamp-to-relative-time
var units = {
  year  : 24 * 60 * 60 * 1000 * 365,
  month : 24 * 60 * 60 * 1000 * 365/12,
  day   : 24 * 60 * 60 * 1000,
  hour  : 60 * 60 * 1000,
  minute: 60 * 1000,
  second: 1000
}

var rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })

var getRelativeTime = (d1, d2 = new Date()) => {
  var elapsed = d1 - d2

  // "Math.abs" accounts for both "past" & "future" scenarios
  for (var u in units) 
    if (Math.abs(elapsed) > units[u] || u == 'second') 
      return rtf.format(Math.round(elapsed/units[u]), u)
}
// endregion

let speech = new SpeechSynthesisUtterance();
speech.lang = 'en';
speech.volume = 1;
speech.rate = 1;
speech.pitch = 1;

function dismissAlert(alertId) {
  console.log('dismissing alert '+alertId);
}

function addAlertElem(alert) {
  let alertElem = document.createElement('article');
  alertElem.classList.add('message');
  alertElem.classList.add('is-danger');
  alertElem.id = alert.id;
  
  // Header
  let alertTitle = document.createElement('div');
  alertTitle.classList.add('message-header');
  let alertTitleTxt = document.createElement('p');
  alertTitleTxt.innerText = `${alert.labels.alertname} - ${alert.labels.name}`;
  let alertDismissBtn = document.createElement('button');
  alertDismissBtn.classList.add('delete');
  alertDismissBtn.onclick = () => dismissAlert(alert.id);
  alertTitle.appendChild(alertTitleTxt);
  alertTitle.appendChild(alertDismissBtn);
  alertElem.appendChild(alertTitle);

  // body
  let alertBody = document.createElement('div');
  alertBody.classList.add('message-body');
  // Get relative time
  let alertTime = document.createElement('p');
  alertTime.innerText = getRelativeTime(new Date(alert.activeAt));
  alertTime.classList.add('has-text-right')
  alertBody.appendChild(alertTime);
  let alertText = document.createElement('p');
  alertText.innerText = `${alert.labels.name} has a ${alert.labels.alertname} alert on ${alert.labels.instance} at `;
  alertBody.appendChild(alertText);
  alertElem.appendChild(alertBody);

  document.getElementById('alerts').appendChild(alertElem);
  speech.text = alertText.innerText;
  console.log(speech)
  window.speechSynthesis.speak(speech);
}

async function fetchAlerts() {
  const rawRes = await fetch(`http://${window.location.host}/api/v1/alerts`);
  const body = await rawRes.json();
  console.log(body.alerts);
  body.alerts.sort((a1, a2) => new Date(a2.activeAt).getTime() - new Date(a1.activeAt).getTime()).forEach(alert => addAlertElem(alert))
}

window.onload = () => {
  fetchAlerts();
  setInterval(fetchAlerts, 5000);
  let endpointElem = document.getElementById('endpoint');
  endpointElem.innerText = window.location.host;
}

window.speechSynthesis.onvoiceschanged = () => {
  // Get List of Voices
  voices = window.speechSynthesis.getVoices();

  // Initially set the First Voice in the Array.
  speech.voice = voices[0];

  // Set the Voice Select List. (Set the Index as the value, which we'll use later when the user updates the Voice using the Select Menu.)
  let voiceSelect = document.querySelector("#voices");
  voices.forEach((voice, i) => (voiceSelect.options[i] = new Option(voice.name, i)));
};