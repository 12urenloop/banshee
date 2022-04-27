let url_to_fetch = ''

function fetchAlerts() {
  const rawRes = await fetch(`http://${}`)
}

window.onload = () => {
  setInterval(fetchAlerts, 5000);
}