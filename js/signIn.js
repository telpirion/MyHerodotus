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
        const errorCode = error.code;
        const errorMessage = error.message;
        const email = error.customData.email;
        const credential = GoogleAuthProvider.credentialFromError(error);

        console.log(`
Error code: ${errorCode}
Error message: ${errorMessage}
Email: ${email}
`)
      });
}

document.getElementById("signIn").addEventListener("click", signIn)

export { signIn };