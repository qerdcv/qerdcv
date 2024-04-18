document.querySelectorAll(".accordion-button").forEach((b) =>
    b.addEventListener("click", (e) => {
        e.target.classList.toggle("show");
        const collapseTarget = document.getElementById(
            e.target.dataset.collapseTarget
        );

        collapseTarget.classList.toggle("show");
        if (collapseTarget.classList.contains("show")) {
            collapseTarget.style.maxHeight = collapseTarget.scrollHeight + "px";
            return;
        }

        collapseTarget.style.maxHeight = "0";
    })
);
