document.addEventListener("DOMContentLoaded", function() {
    let activeBox = null;

    // Function to determine if the device is mobile
    function isMobileDevice() {
    return /Mobi|Android/i.test(navigator.userAgent);
    }

    // Apply readonly attribute to input fields on mobile devices
    if (isMobileDevice()) {
        document.querySelectorAll('.box').forEach(box => {
            box.setAttribute('readonly', true);
        });
    } else {
        document.querySelectorAll('.box').forEach(box => {
            box.removeAttribute('readonly');
        });
    }

    // Focus the first box by default
    const firstBox = document.querySelector('.box');
    if (firstBox) {
        firstBox.focus();
        activeBox = firstBox;
    }

    // Track the last focused input box
    document.querySelectorAll('.box').forEach(box => {
        box.addEventListener('focus', (event) => {
            activeBox = event.target;
        });
    });

    // Function to handle key press
    function handleKeyPress(event) {
        event.preventDefault();  // Prevent the button from stealing focus
        const character = event.target.getAttribute('character').toUpperCase();
        if (activeBox) {
            activeBox.value = character;
            // Optionally, move focus to the next box
            const nextBox = activeBox.nextElementSibling;
            if (nextBox && nextBox.classList.contains('box')) {
                nextBox.disabled = false;
                nextBox.focus();
            }
        }
        sendGridState();
    }

    // Function to handle delete soft-key press
    function handleDeleteKeyPress(event) {
        event.preventDefault();  // Prevent default behavior
        if (activeBox) {
            activeBox.value = '';  // Clear the contents of the active box
            const previousBox = activeBox.previousElementSibling;
            if (previousBox && previousBox.classList.contains('box')) {
                previousBox.focus();
            }
            sendGridState();
        }
    }

    // Function to handle clear soft-key press
    function handleClearKeyPress(event) {
        event.preventDefault();  // Prevent default behavior
        document.querySelectorAll('.box').forEach(box => {
            box.value = '';
        });
        activeBox = firstBox;
        sendGridState();
    }

    // Add event listeners to each keyboard button
    const keyboardKeys = document.querySelectorAll('.keyboardKey:not(.keyboardKey[character="←"]):not(.keyboardKey[character="clear"])');
    keyboardKeys.forEach(key => {
        key.addEventListener('click', handleKeyPress);
    });

    // Add event listener to the delete soft-key button
    const deleteKey = document.querySelector('.keyboardKey[character="←"]');
    deleteKey.addEventListener('click', handleDeleteKeyPress);

    // Add event listener to the clear soft-key button
    const clearKey = document.querySelector('.keyboardKey[character="clear"]');
    clearKey.addEventListener('click', handleClearKeyPress);

    // Track the last focused input box
    document.querySelectorAll('.box').forEach(box => {
        box.addEventListener('focus', (event) => {
            activeBox = event.target;
        });

        // Add event listener for backspace key press (on physical keyboard)
        box.addEventListener('keydown', (event) => {
            if (event.key === 'Backspace' && box.value === '') {
                event.preventDefault();  // Prevent default backspace behavior
                const previousBox = box.previousElementSibling;
                if (previousBox && previousBox.classList.contains('box')) {
                    previousBox.focus();
                }
                sendGridState();
            }
        });
    });

    function limitInputToSingleChar(event) {
        const input = event.target;
        if (input.value.length > 1) {
            input.value = input.value.charAt(0);
        }
        // Enable the next box if this one is filled
        if (input.value.length === 1) {
            enableNextBox(input);
        }
        sendGridState();
    }

    function enableNextBox(currentBox) {
        const boxes = Array.from(document.querySelectorAll('.box'));
        const currentIndex = boxes.indexOf(currentBox);
        if (currentIndex < boxes.length - 1) {
            boxes[currentIndex + 1].disabled = false;
            boxes[currentIndex + 1].focus();
        }
    }

    async function sendGridState() {
        const boxes = document.querySelectorAll('.box');
        const gridData = Array.from(boxes).map(box => ({
            character: box.value.toUpperCase(),
            color: box.className.split(' ').find(cls => ['grey', 'yellow'].includes(cls)) || ''
        }));
    
        try {
            const response = await fetch('/letters', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(gridData)
            });
            const data = await response.json();
            document.getElementById('resultsCount').innerText = data.resultCount;
            document.getElementById('results').innerText = data.result;
            if (data.error) {
                throw document.getElementById('resultsCount').innerText = 0,
                document.getElementById('results').innerText = data.error;
            }
        } catch (error) {
            console.error('error:', error);   
        }
    }

    const boxes = document.querySelectorAll('.box');
    boxes.forEach(box => {
        box.addEventListener('input', limitInputToSingleChar);
    });
});

// Function to apply the theme
function applyTheme(theme) {
    var darkThemeToggleButton = document.getElementById("darkThemeToggleButton");
    var headerNav = document.getElementById("headerNav");
    var title = document.getElementById("title");
    var body = document.getElementById("body");
    var pageFrame = document.getElementById("pageFrame");
    var boxes = document.getElementsByClassName("box");
    var answerbox = document.getElementById("answerbox");

    if (theme === "dark") {
        // Apply dark theme styles
        // Toggle the button
        darkThemeToggleButton.className = "fa fa-toggle-on fa-2x";
        darkThemeToggleButton.style.backgroundColor = "rgb(35, 35, 35)";
        darkThemeToggleButton.style.color = "white";
        // Header colours
        headerNav.style.backgroundColor = "rgb(35, 35, 35)";
        headerNav.style.color = "white";
        // Change title colours
        title.style.backgroundColor = "rgb(35, 35, 35)";
        title.style.color = "white";
        // Change background colour
        body.style.backgroundColor = "rgb(35, 35, 35)";
        // Change instructions colour
        pageFrame.style.color = "white";
        // Change answer box colours
        answerbox.style.backgroundColor = "rgb(35, 35, 35)";
        answerbox.style.color = "white";
    } else {
        // Apply light theme styles
        // Toggle the button
        darkThemeToggleButton.className = "fa fa-toggle-off fa-2x";
        darkThemeToggleButton.style.backgroundColor = "white";
        darkThemeToggleButton.style.color = "black";
        //Header colours
        headerNav.style.backgroundColor = "white";
        headerNav.style.color = "black";
        //Change title colours
        title.style.backgroundColor = "white";
        title.style.color = "black";
        //Change background colour
        body.style.backgroundColor = "white";
        //Change instructions colour
        pageFrame.style.color = "black";
        //Change answer box colours
        answerbox.style.backgroundColor = "white";
        answerbox.style.color = "black";
    }
}

// Function to toggle the theme and save preference
function toggleTheme() {
    var currentTheme = localStorage.getItem("theme");

    if (currentTheme === "dark") {
        // Switch to light theme and save preference
        applyTheme("light");
        localStorage.setItem("theme", "light");
    } else {
        // Switch to dark theme and save preference
        applyTheme("dark");
        localStorage.setItem("theme", "dark");
    }
}

// Check localStorage for the theme preference when the page loads
window.onload = function() {
    var savedTheme = localStorage.getItem("theme") || "light"; // Default to light theme
    applyTheme(savedTheme);
};
