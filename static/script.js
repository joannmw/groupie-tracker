function openTab(evt, tabName) {
    var i, tabContent, tabLinks;
    tabContent = document.getElementsByClassName("tab-content");
    for (i = 0; i < tabContent.length; i++) {
        tabContent[i].style.display = "none";
    }
    tabLinks = document.getElementsByClassName("tab");
    for (i = 0; i < tabLinks.length; i++) {
        tabLinks[i].className = tabLinks[i].className.replace(" active", "");
    }
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}

const cursor = document.createElement('div');
cursor.classList.add('custom-cursor');
document.body.appendChild(cursor);

document.addEventListener('mousemove', (e) => {
    cursor.style.left = e.clientX + 'px';
    cursor.style.top = e.clientY + 'px';

    // Check if the element under cursor is clickable
    const elementUnderCursor = document.elementFromPoint(e.clientX, e.clientY);
    if (elementUnderCursor) {
        const isClickable = (
            elementUnderCursor.matches('a, button, [role="button"], input, .clickable') ||
            elementUnderCursor.closest('a, button, [role="button"], input, .clickable')
        );

        // Set opacity based on whether element is clickable
        cursor.style.opacity = isClickable ? '0' : '1';
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const tooltip = document.getElementById('tooltip');
    const buttons = document.querySelectorAll('.search-button, .details-button, .home-link');

    buttons.forEach(button => {
        button.addEventListener('mouseenter', (event) => {
            tooltip.textContent = button.getAttribute('data-tooltip');
            tooltip.style.display = 'block';
            // Get the position of the button and set the tooltip below it
            const rect = event.target.getBoundingClientRect();
            tooltip.style.left = `${rect.left + window.scrollX + rect.width / 2 - tooltip.offsetWidth / 2}px`;
            tooltip.style.top = `${rect.bottom + window.scrollY + 5}px`;
        });

        button.addEventListener('mouseleave', () => {
            tooltip.style.display = 'none';
        });
    });
});