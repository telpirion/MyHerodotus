import { getAuth } from 'https://www.gstatic.com/firebasejs/10.14.1/firebase-auth.js'

function processForm(e) {
    //e.preventDefault();

    // Emit 'msg' event for bot progress bar
    const event = new Event("msg");
    document.dispatchEvent(event);

    // Get the user email address
    const auth = getAuth();
    const user = auth.currentUser;
    
    // If the user hasn't signed in, redirect them to the sign-in page.
    if (!user) {
        window.location = `/?status=unauthorized`;
    }

    // Collect data
    const message = document.getElementById("userMsg").value;
    const selection = document.getElementById("modelSelect");
    const model = selection.options[selection.selectedIndex].text;

    const payload = {
        model,
        message,
    };

    fetch(`/home?user=${user.email}`, {
        method: "POST",
        body: JSON.stringify(payload),
        headers: {
            "Content-Type": "application/json",
        }
    })
    .then(response => {
        return response.json();
    })
    .then(data => {
        console.log(data);
        const response = data.Message.Message;
        const botMessage = document.querySelector(".message-actual")
        botMessage.textContent = response;
        botMessage.dataset.document = data.Message.DocumentID;

        // Toggle visibility
        const evt = new Event("msg");
        document.dispatchEvent(evt);
    })
    return true;
}

window.addEventListener("load", function() {
    const send = document.getElementById("send");
    if (send.attachEvent) {
        send.attachEvent("click", processForm);
    } else {
        send.addEventListener("click", processForm);
    }
});