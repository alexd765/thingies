const bucketInput = 'thingies-input'
const bucketOutput = 'thingies-output'
const region = 'eu-west-1'
const identityPool = 'eu-west-1:c4de6544-663a-40bd-a017-2e0924f55ca1'

AWS.config.update({
    region: region,
    credentials: new AWS.CognitoIdentityCredentials({
        IdentityPoolId: identityPool
    })
})

const s3 = new AWS.S3({
    apiVersion: '2006-03-01',
})

function upload() {
    const files = document.getElementById("select-files").files
    if (!files.length) {
        msg('please select an image')
        return
    }
    const file = files[0]
    const token = generateToken(file.name)

    msg('uploading...')
    s3.upload({
        Key: token,
        Bucket: bucketInput,
        Body: file,
    }, convert)
}

function convert(err, data) {
    if (err) {
        msg('upload failed: ' + err.message)
        return
    }
    msg('upload sucessful. converting ...')

    var req = new XMLHttpRequest()
    req.open("POST", "https://vvk4r5qg29.execute-api.eu-west-1.amazonaws.com/prod")
    req.setRequestHeader("Content-Type", "application/json")
    req.onreadystatechange = download
    req.send(JSON.stringify({ version: 1, "input-token": data.key }))

    displayImageInput(data.key)
}

function download() {
    if (this.readyState != XMLHttpRequest.DONE) {
        return
    }
    const resp = JSON.parse(this.responseText)
    if (resp.version != 3 || resp["output-token"] == "") {
        msg('convert unsuccessful')
        return
    }
    console.log("resp", resp)
    msg('convert sucessful. downloading...')
    displayImageOutput(resp["output-token"])
}

function generateToken(filename) {
    const ext = filename.substr(filename.lastIndexOf('.'))
    return rand() + ext
}

function msg(text) {
    document.getElementById("message").innerHTML = text
}

function displayImageInput(token) {
    const url = s3.getSignedUrl('getObject', {
        Key: token,
        Bucket: bucketInput
    })
    document.getElementById("input-image").src = url
}

function displayImageOutput(token) {
    console.log("token: ", token)
    const url = s3.getSignedUrl('getObject', {
        Key: token,
        Bucket: bucketOutput
    })
    document.getElementById("output-image").src = url
}

function rand() {
    const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var text = ""
    for (var i = 0; i < 16; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length))
    }
    return text
}