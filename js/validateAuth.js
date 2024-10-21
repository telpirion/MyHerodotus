import { getAuth, onAuthStateChanged } from 'https://www.gstatic.com/firebasejs/10.14.1/firebase-auth.js'

window.addEventListener("load", function () {
    const auth = getAuth();
    const user = auth.currentUser;
    
    onAuthStateChanged(auth, (user) => {
        if (user) {
          const uid = user.uid;
        } else {
            window.location = `/?message=sign-in-required`;
        }
    });
});
