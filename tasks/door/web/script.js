function getBaseXY() {
	const lockRect = lock.getBoundingClientRect();
	const baseX = (lockRect.left + lockRect.right) / 2;

	const handleRect = handle.getBoundingClientRect();
	const baseY = (handleRect.top + handleRect.bottom) / 2;

	return [baseX, baseY];
}

let isMovingHandle = false;
let handleTouchId = null;
let handleRadius, baseX, baseY;
let angle = 0;
let story = [];

function updateHandle() {
	handle.style.transform = `perspective(10vw) rotate(${Math.floor(angle / Math.PI * 180)}deg) rotateY(40deg)`;
}

function onStart(touchId, target, x, y) {
	if (target === door) {
		story.push(["door"]);
		fetch("/enter", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(story)
		}).then(res => res.text()).then(alert);
	}
	if (isMovingHandle || target !== handle) {
		return;
	}
	handleTouchId = touchId;
	[baseX, baseY] = getBaseXY();
	const dx = baseX - x;
	const dy = baseY - y;
	story.push(["start", touchId, target.id, dx, dy]);
	const radius = Math.sqrt(dx * dx + dy * dy);
	if (radius > handle.getBoundingClientRect().width / 2) {
		isMovingHandle = true;
		handleRadius = radius;
	}
}

function onMove(touchId, x, y) {
	if (!isMovingHandle || handleTouchId !== touchId) {
		return;
	}
	const dx = x - baseX;
	const dy = y - baseY;
	story.push(["move", touchId, dx, dy]);
	const radius = Math.sqrt(dx * dx + dy * dy);
	if (dx < 0 || Math.abs(radius - handleRadius) > 8) {
		angle = 0;
		isMovingHandle = false;
	} else {
		angle = Math.min(Math.max(-0.1, Math.atan2(dy, dx)), 1.2);
	}
	updateHandle();
}

function onStop(touchId) {
	if (handleTouchId !== touchId) {
		return;
	}
	story.push(["stop", touchId]);
	angle = 0;
	isMovingHandle = false;
	updateHandle();
}

function on(name, handler) {
	window.addEventListener(name, handler);
	// window.addEventListener(name, e => {
	// 	e.preventDefault();
	// 	handler(e);
	// }, {passive: false});
}

on("mousedown", e => {
	onStart("mouse", e.target, e.clientX, e.clientY);
});
on("touchstart", e => {
	for (const touch of e.changedTouches) {
		onStart(touch.identifier, touch.target, touch.clientX, touch.clientY);
	}
});

on("mousemove", e => {
	onMove("mouse", e.clientX, e.clientY);
});
on("touchmove", e => {
	for (const touch of e.changedTouches) {
		onMove(touch.identifier, touch.clientX, touch.clientY);
	}
});

on("mouseup", e => {
	onStop("mouse");
});
on("touchend", e => {
	for (const touch of e.changedTouches) {
		onStop(touch.identifier);
	}
});
