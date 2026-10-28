fetch('/api/items')
  .then((r) => r.json())
  .then((items) => {
    document.getElementById('items').textContent = JSON.stringify(items, null, 2);
  })
  .catch((err) => {
    document.getElementById('items').textContent = 'Ошибка: ' + err.message;
  });