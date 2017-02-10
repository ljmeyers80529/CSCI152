var entry = document.querySelector('#enterUsername');
var output = document.querySelector('#usedmsg');

function inputProcess()
{
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/username/check');
    xhr.send(entry.value);
    xhr.addEventListener('readystatechange', function(){
        if (xhr.readyState === 4 && xhr.status === 200) {
            var taken = xhr.responseText;
            if (taken == 'true') {
                output.textContent = 'Sorry, user name already taken!';
            } else {
                output.textContent = '';
            }
        }
    });
}

entry.addEventListener('input', inputProcess);
