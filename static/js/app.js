document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll(".reactions").forEach(function (group) {
        var messageId = group.dataset.messageId;

        group.querySelectorAll(".reaction-btn").forEach(function (button) {
            button.addEventListener("click", function () {
                sendReaction(group, messageId, button.dataset.type);
            });
        });
    });
});

function sendReaction(group, messageId, type) {
    fetch("/api/messages/" + messageId + "/reactions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type: type })
    })
        .then(function (response) {
            if (!response.ok) {
                return null;
            }
            return response.json();
        })
        .then(function (data) {
            if (!data) {
                return;
            }

            group.querySelector(".count-like").textContent = data.likes;
            group.querySelector(".count-dislike").textContent = data.dislikes;
            group.querySelector(".count-score").textContent = data.score;

            group.querySelectorAll(".reaction-btn").forEach(function (button) {
                button.classList.remove("active");
            });

            if (data.reaction_utilisateur === "like") {
                group.querySelector(".reaction-btn.like").classList.add("active");
            } else if (data.reaction_utilisateur === "dislike") {
                group.querySelector(".reaction-btn.dislike").classList.add("active");
            }
        })
        .catch(function () {});
}
