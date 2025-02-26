window.onload = function() {
    const flag1 = { x: 30, y: 30, width: 40, height: 40 };
    const flag2 = { x: 900, y: 400, width: 40, height: 40 };

    
    const startX = flag1.x + 2; // Левый край
    const startY = flag1.y + flag1.height; // Нижний край
    const endX = flag2.x + 2; // Левый край
    const endY = flag2.y + flag2.height; // Нижний край

    // Создаем путь для извилистой линии с двумя изгибами
    const controlX1 = startX + (endX - startX) / 3; 
    const controlY1 = startY + 400; 

    const controlX2 = startX + (endX - startX) * 0.7; 
    const controlY2 = startY - 200; 

    const d = `M ${startX} ${startY} C ${controlX1} ${controlY1}, ${controlX2} ${controlY2}, ${endX} ${endY}`;
    document.getElementById('line').setAttribute('d', d);
    document.getElementById('dashed-path').setAttribute('d', d);
};