function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(
        function () {
            showTooltip();
        },
        function (err) {
            console.error("Could not copy text: ", err);
        },
    );
}

function showTooltip() {
    var tooltip = document.getElementById("tooltip");
    tooltip.classList.remove("d-none");
    setTimeout(function () {
        tooltip.classList.add("d-none");
    }, 5000);
}

function showAlert(alertId, message) {
    var alertElement = document.getElementById(alertId);
    if (message) {
        alertElement.innerHTML = message;
    }
    alertElement.classList.remove("d-none");
}

function hideAlert(alertId) {
    var alertElement = document.getElementById(alertId);
    alertElement.classList.add("d-none");
}

document.body.addEventListener(
    "htmx:afterRequest",
    function (event) {
        hideAlert("result");
        hideAlert("too-many-request");
        hideAlert("not-found");

        if (event.detail.xhr.status === 404) {
            showAlert("not-found", event.detail.xhr.responseText);
        } else if (event.detail.xhr.status === 429) {
            var retryAfter = event.detail.xhr.getResponseHeader("Retry-After");
            var message = "";
            if (retryAfter) {
                message = ` Try again after <strong>${retryAfter}</strong> seconds.`;
            }
            showAlert(
                "too-many-request",
                "<div>" +event.detail.xhr.responseText + message + "</div>",
            );
            setTimeout(function () {
                hideAlert("too-many-request");
            }, retryAfter * 1000);
        } else if (event.detail.xhr.status === 200) {
            showAlert("result");
        } else if (event.detail.xhr.status === 0) {
            showAlert(
                "not-found",
                "Endpoint could not be reached. Please try again later.",
            );
        }
    },
);

document.body.addEventListener("htmx:afterSwap", function (event) {
    if (event.detail.target.id === "result") {
        showAlert("result");
    }
});

document.body.addEventListener("htmx:load", function (event) {
    var domainUrl = window.location.origin;
    const url = new URL(domainUrl);
    const domain = url.hostname;

    var banner = document.getElementById("banner");
    banner.innerHTML = `<h2 class="text-light">${domain}</h2><p class="text-light">Make your links shorter</p>`;
    b
});