window.addEventListener("load", function () {

    document.addEventListener("msg", () => {
        document.querySelector("progress").classList.toggle("is-hidden");
        document.querySelector(".message-actual").classList.toggle("is-hidden");
    });

    document.querySelectorAll(".thumb").forEach(e => {
        e.addEventListener("click", e => {
            const rating = e.currentTarget.getAttribute("id");
            const botMessage = this.document.querySelector(".message-actual")
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
            })
        });
    });
});