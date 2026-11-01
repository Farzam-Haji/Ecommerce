// const root = document.documentElement;
// const themeSelect = document.getElementById("theme-select");
// const systemTheme = window.matchMedia("(prefers-color-scheme: dark)");


// function applyTheme(theme) {
//     if (theme === "system") {
//         root.dataset.theme = systemTheme.matches ? "dark" : "light";
//         return;
//     }

//     root.dataset.theme = theme;
// }


// themeSelect.addEventListener("change", function () {
//     const selectedTheme = themeSelect.value;

//     localStorage.setItem("theme", selectedTheme);
//     applyTheme(selectedTheme);
// });


// systemTheme.addEventListener("change", function () {
//     const savedTheme = localStorage.getItem("theme");

//     if (savedTheme === "system" || savedTheme === null) {
//         applyTheme("system");
//     }
// });


// // initial page load
// const savedTheme = localStorage.getItem("theme");
// const initialTheme = savedTheme || "system";

// themeSelect.value = initialTheme;
// applyTheme(initialTheme);



const root = document.documentElement;

const themeMenuButton = document.getElementById("theme-menu-button");
const themeMenu = document.getElementById("theme-menu");

const themeOptions = document.querySelectorAll("[data-theme-option]");

const systemTheme = window.matchMedia("(prefers-color-scheme: dark)");


// Apply the selected theme
function applyTheme(theme) {
    if (theme === "system") {
        root.dataset.theme = systemTheme.matches ? "dark" : "light";
    } else {
        root.dataset.theme = theme;
    }

    updateThemeIcon(theme);
}


// Update the icon in the menu trigger
function updateThemeIcon(theme) {
    document.getElementById("theme-icon-system").classList.toggle("hidden", theme !== "system");
    document.getElementById("theme-icon-light").classList.toggle("hidden", theme !== "light");
    document.getElementById("theme-icon-dark").classList.toggle("hidden", theme !== "dark");
}


// Open and close the dropdown
themeMenuButton.addEventListener("click", function () {
    const isOpen = !themeMenu.classList.contains("hidden");

    themeMenu.classList.toggle("hidden");

    themeMenuButton.setAttribute("aria-expanded", String(!isOpen));
});


// Handle selecting a theme
themeOptions.forEach(function (option) {
    option.addEventListener("click", function () {
        const selectedTheme = option.dataset.themeOption;

        localStorage.setItem("theme", selectedTheme);

        applyTheme(selectedTheme);

        themeMenu.classList.add("hidden");
        themeMenuButton.setAttribute("aria-expanded", "false");
    });
});


// Close menu when clicking outside
document.addEventListener("click", function (event) {
    if (!themeMenuButton.parentElement.contains(event.target)) {
        themeMenu.classList.add("hidden");
        themeMenuButton.setAttribute("aria-expanded", "false");
    }
});


// React to operating system theme changes
systemTheme.addEventListener("change", function () {
    if (localStorage.getItem("theme") === "system") {
        applyTheme("system");
    }
});


// Initialize theme
const savedTheme = localStorage.getItem("theme");
const initialTheme = savedTheme || "system";

applyTheme(initialTheme);
