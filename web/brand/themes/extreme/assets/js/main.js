document.addEventListener("DOMContentLoaded", function () {
    // Select all dynamic word elements
    const dynamicWordElements = document.querySelectorAll(".dynamic-word");

    // Function to initialize word cycling on a single element
    const initWordCycling = (element) => {
        const fadeDuration = 500; // Duration for fade-in and fade-out animations in milliseconds
        const wordCycleInterval = 3000; // Interval between word changes in milliseconds
        const words = element.getAttribute('data-words').split(',');
        let currentWordIndex = 0;

        // Set interval to cycle through words
        setInterval(() => {
            element.classList.add('fade-out-up'); // Start the fade-out animation

            setTimeout(() => {
                currentWordIndex = (currentWordIndex + 1) % words.length;
                element.textContent = words[currentWordIndex];
                element.classList.remove('fade-out-up'); // Remove the fade-out class
                element.classList.add('fade-in-up');     // Add fade-in class

                // Remove fade-in class after animation completes to reset the state
                setTimeout(() => {
                    element.classList.remove('fade-in-up');
                }, fadeDuration); // Matches the animation duration
            }, fadeDuration); // Matches the fade-out animation duration
        }, wordCycleInterval); // Change word at set interval
    };

    // Initialize word cycling for each element found
    dynamicWordElements.forEach(element => {
        initWordCycling(element);
    });
});