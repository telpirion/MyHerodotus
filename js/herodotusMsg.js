window.addEventListener("load", function () {

    document.addEventListener("msg", () => {
        document.querySelector("progress").classList.toggle("is-hidden");
        document.querySelector(".message-actual").classList.toggle("is-hidden");
    });

    document.querySelectorAll(".thumb").forEach(e => {
        e.addEventListener("click", e => {
            const rating = e.target.getAttribute("id");
            const response = document.getElementById("messageBody").innerText;
            const payload = {
                rating,
                response,
            }

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
            })
        });
    });
});