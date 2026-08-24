const defaults = {
  coupleNames: "Lin & Kai",
  date: "2026-08-24",
  venue: "Jasmine Hall",
  welcomeText: "Welcome to our celebration",
  heroImage: ""
};

function applyRecord(record) {
  const data = Object.assign({}, defaults, record || {});
  document.querySelector("#couple-names").textContent = data.coupleNames;
  document.querySelector("#wedding-date").textContent = data.date;
  document.querySelector("#wedding-venue").textContent = data.venue;
  document.querySelector("#welcome-text").textContent = data.welcomeText;
  const shell = document.querySelector("#welcome");
  shell.style.backgroundImage = data.heroImage ? `url('${data.heroImage}')` : "linear-gradient(135deg, #273a35, #b48b65)";
}

function startFullscreen() {
  const shell = document.querySelector("#welcome");
  if (document.fullscreenElement) return;
  if (shell.requestFullscreen) shell.requestFullscreen();
}

document.querySelector("#start-button").addEventListener("click", startFullscreen);
applyRecord(window.weddingRecord);
