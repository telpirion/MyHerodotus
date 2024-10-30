import { getAuth, signInWithPopup, GoogleAuthProvider } from 'https://www.gstatic.com/firebasejs/10.14.1/firebase-auth.js'

const signIn = () => {
    const auth = getAuth();
    const provider = new GoogleAuthProvider();
    signInWithPopup(auth, provider)
      .then((result) => {
        // This gives you a Google Access Token. You can use it to access the Google API.
        const credential = GoogleAuthProvider.credentialFromResult(result);
        const token = credential.accessToken;
    
        const user = result.user;

        window.location = `/home?user=${user.email}`;
      }).catch((error) => {
        const code = error.code;
        const message = error.message;
        const email = error.customData.email;
        const credential = GoogleAuthProvider.credentialFromError(error);

        console.log(`
Error code: ${code}
Error message: ${message}
Email: ${email}
`)
        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/logClientError", true)
        xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        xhr.send(JSON.stringify( 
          {
            code,
            email,
            message,
            credential,
          }
        ));
        window.location = "/error"
      });
}

document.getElementById("signIn").addEventListener("click", signIn)

export { signIn };