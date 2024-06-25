let n = 0;
document.addEventListener('mousemove', function (e) {
    if (n < 100) {
        createParticle(e.clientX, e.clientY);
        n += 1;
    }
    console.log(n);
});

function createParticle(x, y) {
    const particle = document.createElement('div');
    particle.classList.add('particle');
    particle.style.left = `${x}px`;
    particle.style.top = `${y}px`;

    document.body.appendChild(particle);

    setTimeout(() => {
        particle.style.transform = `translate(${random(-50, 50)}px, ${random(-50, 50)}px)`;
        particle.style.opacity = '0';
    }, 0);

    setTimeout(() => {
        n -= 1;
        particle.remove();
    }, random(10,100));
}

function random(min, max) {
    return Math.random() * (max - min) + min;
}
