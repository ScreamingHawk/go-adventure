const addMessage = (msg, isUser) => {
  const msgList = document.getElementById("msgList");
  const p = document.createElement("p");
  p.className = isUser ? "msg-user" : "msg-ai";
  p.textContent = msg;
  msgList.appendChild(p);
};

// Generate a random ID
const sessionId = Math.random().toString(36).substring(2, 9);

const useResponse = async (response) => {
  if (!response.ok) {
    throw new Error("Request failed.");
  }
  const resp = await response.json();

  addMessage(resp.plot, false);

  const btnBox = document.getElementById("btnBox");
  if (resp.choices && resp.choices.length > 0) {
    resp.choices.forEach((choice) => {
      const btn = document.createElement("button");
      btn.className = "choiceBtn";
      btn.textContent = choice;
      btn.addEventListener("click", updateClick);
      btnBox.appendChild(btn);
    });
  } else {
    // End of story
    const endP = document.createElement("p");
    endP.textContent = "The end.";
    btnBox.appendChild(endP);
  }
};

const startClick = (event) => {
  // Clear buttons and message list
  document.getElementById("msgList").innerHTML = "";
  event.target.parentElement.innerHTML = "";

  fetch(`/api/narrate/${sessionId}`, {
    method: "GET",
  }).then(useResponse);
};

const updateClick = (event) => {
  const choice = event.target.textContent;
  addMessage(choice, true);

  // Clear buttons
  event.target.parentElement.innerHTML = "";

  fetch(`/api/narrate/${sessionId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ choice }),
  }).then(useResponse);
};

document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("startBtn").addEventListener("click", startClick);
});
