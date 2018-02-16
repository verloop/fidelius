# Fidelius

Read [this blog post](https://medium.com/verloop-engineering/fidelius-e4b2d8b6b1df) for more details.

## Usage

If you want to generate a temporary access token:

    $ fidelius --gh-integration-id=1234 --gh-installation-id=5678 --gh-private-key='/root/fidelius-charm.2018-02-12.private-key.pem'
    # will print the token
    v1.437ad23ce1a4a71b1f3222bc198de653619a9570

If you want to generate a proper `.gitconfig` file:

    $ fidelius --gh-integration-id=1234 --gh-installation-id=5678 --gh-private-key='/root/fidelius-charm.2018-02-12.private-key.pem' --git-config-out='/root/.gitconfig'

or if you want to modify an existing `.gitconfig`:

    $ export GH_TOKEN=$(fidelius --gh-integration-id=1234 --gh-installation-id=5678 --gh-private-key='/root/fidelius-charm.2018-02-12.private-key.pem')
    $ git config --global url.https://x-access-token:$GH_TOKEN@github.com/.insteadOf https://github.com/

## Name

[Fidelius Charm](http://harrypotter.wikia.com/wiki/Fidelius_Charm)

## License

Released under Unlicense. Check the `LICENSE` for more details.