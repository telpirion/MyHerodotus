import { getAuth,  onAuthStateChanged } from 'https://www.gstatic.com/firebasejs/10.14.1/firebase-auth.js'

window.addEventListener("load", function () {
    const send = document.querySelector('#send');
    const auth = getAuth();
    let userEmail = ""

    if (send.attachEvent) {
        send.attachEvent("click", processForm);
    } else {
        send.addEventListener("click", processForm);
    }
    
    onAuthStateChanged(auth, (user) => {
        if (user) {
          userEmail = user.email;
        } else {
            window.location = `/?message=sign-in-required`;
        }
    });

    document.addEventListener("msg", toggleProgressBar);

    document.querySelectorAll(".thumb").forEach(e => {
        e.addEventListener("click", processRating);
    });
});

function toggleProgressBar() {
    document.querySelector("progress").classList.toggle("is-hidden");
    document.querySelector(".message-actual").classList.toggle("is-hidden");
}

function processRating(e) {
    const toast = window.document.querySelector('.notification');
    const rating = e.currentTarget.getAttribute("id");
    const botMessage = window.document.querySelector(".message-actual")
    const document = botMessage.dataset.document;
    const response = botMessage.textContent;
    const payload = {
        rating,
        response,
        document,
    }

    fetch(`/rateResponse`, {
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
        toast.classList.toggle("toast");
        toast.classList.toggle("toast-hide");
        window.setTimeout(()=>{
            toast.classList.toggle("toast-hide");
            toast.classList.toggle("toast");
        }, 5000);
    });
}

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
