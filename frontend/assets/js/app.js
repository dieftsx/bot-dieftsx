const API_URL = "http://127.0.0.1:8080/api/data";

let lastFollower = null;
let lastSub = null;
let lastBitsUser = null;
let lastHypeLevel = null;
const alertQueue = [];
let isAnimating = false;

async function fetchStreamData() {
  try {
    const response = await fetch(API_URL);
    const data = await response.json();
    updateGeneralUI(data);
    checkAlerts(data);
  } catch (error) {
    console.error("Erro no backend:", error);
  }
}

function updateGeneralUI(data) {
  const el = (id) => document.getElementById(id);

  if (el("status-text"))
    el("status-text").innerHTML =
      `<span class="text-green">●</span> ${data.status}`;

  if (el("sys-info")) {
    const sysList = el("sys-info");
    sysList.innerHTML = "";
    for (const [key, value] of Object.entries(data.system_info)) {
      if (sysList.tagName === "UL") {
        sysList.innerHTML += `<li style="margin-bottom:8px;"><strong>${key}:</strong> <span style="color:var(--cyan);">${value}</span></li>`;
      } else {
        sysList.innerHTML += `<div><span style="color:var(--cyan);">${key}</span> <span style="color:var(--text-light);">${value}</span></div>`;
      }
    }
  }

  if (el("dir-path")) el("dir-path").innerText = data.project.directory;
  if (el("git-branch")) el("git-branch").innerText = data.project.branch;
  if (el("git-status")) {
    el("git-status").innerText = data.project.status;
    el("git-status").style.color = data.project.status.includes("Clean")
      ? "var(--green)"
      : "var(--yellow)";
  }

  if (el("stat-duration")) el("stat-duration").innerText = data.stats.duration;
  if (el("stat-lines")) el("stat-lines").innerText = data.stats.lines;
  if (el("stat-commits")) el("stat-commits").innerText = data.stats.commits;
  if (el("stat-messages")) el("stat-messages").innerText = data.stats.messages;
}

function checkAlerts(data) {
  if (!document.getElementById("follow-container")) return;

  if (
    lastFollower !== null &&
    data.latest_follower !== lastFollower &&
    data.latest_follower !== "Aguardando..."
  )
    alertQueue.push({ id: "follow-container", name: data.latest_follower });
  lastFollower = data.latest_follower;

  if (
    lastSub !== null &&
    data.latest_sub !== lastSub &&
    data.latest_sub !== "Aguardando..."
  )
    alertQueue.push({ id: "sub-container", name: data.latest_sub });
  lastSub = data.latest_sub;

  if (
    lastBitsUser !== null &&
    data.latest_bits_user !== lastBitsUser &&
    data.latest_bits_user !== "Aguardando..."
  )
    alertQueue.push({
      id: "bits-container",
      name: data.latest_bits_user,
      amount: data.latest_bits,
    });
  lastBitsUser = data.latest_bits_user;

  if (
    lastHypeLevel !== null &&
    data.hype_train_level !== lastHypeLevel &&
    data.hype_train_level !== "0"
  )
    alertQueue.push({ id: "hype-container", level: data.hype_train_level });
  lastHypeLevel = data.hype_train_level;

  processQueue();
}

function processQueue() {
  if (isAnimating || alertQueue.length === 0) return;

  isAnimating = true;
  const alertData = alertQueue.shift();
  const container = document.getElementById(alertData.id);

  if (alertData.id === "follow-container")
    document.getElementById("alert-follower-name").innerText = alertData.name;
  if (alertData.id === "sub-container")
    document.getElementById("alert-sub-name").innerText = alertData.name;
  if (alertData.id === "bits-container") {
    document.getElementById("alert-bits-name").innerText = alertData.name;
    document.getElementById("alert-bits-amount").innerText = alertData.amount;
  }
  if (alertData.id === "hype-container")
    document.getElementById("alert-hype-level").innerText = alertData.level;

  container.classList.add("show");
  setTimeout(() => {
    container.classList.remove("show");
    setTimeout(() => {
      isAnimating = false;
      processQueue();
    }, 600);
  }, 5000);
}

setInterval(fetchStreamData, 1000);
fetchStreamData();

