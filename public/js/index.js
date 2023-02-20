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

let loadedAlerts = [];
let audio = new Audio('public/assets/siren.mp3');

function dismissAlert(alertId) {
  fetch(`http://${window.location.host}/api/v1/alerts`, {
    method: 'PUT',
    body: alertId,
  })
  removeAlert(alertId);
}

const removeAlert = (alertId) => {
  let alertElem = document.getElementById(alertId);
  // delete alertElem;
  alertElem.remove();
  loadedAlerts = loadedAlerts.filter(a => a !== alertId);
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
  let alertText = document.createElement('p');
  alertText.innerText = `${alert.labels.name} has a ${alert.labels.alertname} alert on ${alert.labels.instance}`;
  alertBody.appendChild(alertText);
  // Get relative time
  let alertTime = document.createElement('p');
  alertTime.innerText = getRelativeTime(new Date(alert.activeAt));
  alertTime.classList.add('has-text-right')
  alertBody.appendChild(alertTime);
  alertElem.appendChild(alertBody);
  
  setInterval(() => {
    alertTime.innerText = getRelativeTime(new Date(alert.activeAt));
  },1000)

  // Play airhorn sound
  audio.play();

  document.getElementById('alerts').appendChild(alertElem);
  loadedAlerts.push(alert.id);
}

async function fetchAlerts() {
  const rawRes = await fetch(`http://${window.location.host}/api/v1/alerts`);
  const body = await rawRes.json();
  body.alerts.filter(a => !loadedAlerts.includes(a.id)).sort((a1, a2) => new Date(a2.activeAt).getTime() - new Date(a1.activeAt).getTime()).forEach(alert => addAlertElem(alert))
  body.dismissedAlerts.filter(id => loadedAlerts.includes(id)).forEach(id => removeAlert(id))
}

window.onload = () => {
  fetchAlerts();
  setInterval(fetchAlerts, 5000);
  let endpointElem = document.getElementById('endpoint');
  endpointElem.innerText = window.location.host;
}
