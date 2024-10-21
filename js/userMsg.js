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
    return true;
}

window.addEventListener("load", function() {
    const form = document.getElementById("userMsg");
    if (form.attachEvent) {
        form.attachEvent("submit", processForm);
    } else {
        form.addEventListener("submit", processForm);
    }
});