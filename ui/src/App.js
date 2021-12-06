import React from 'react';
import './App.css';
import web3 from './web3';

class App extends React.Component {
    raffle = null;

    state = {
        manager: '',
        players: [],
        balance: '',
        value: '',
        message: ''
    };

    async componentDidMount() {
        const rAddress = await fetch('http://localhost:8080/contract-address');
        const address = await rAddress.text();

        const rAbi = await fetch('http://localhost:8080/contract-abi');
        const abi = await rAbi.json();

        this.raffle = new web3.eth.Contract(abi, address);

        const manager = await this.raffle.methods.manager().call();
        const players = await this.raffle.methods.getPlayers().call();
        const balance = await web3.eth.getBalance(this.raffle.options.address);

        this.setState({manager, players, balance});
    }

    onSubmit = async event => {
        event.preventDefault();

        const accounts = await web3.eth.getAccounts();

        this.setState({message: 'Waiting on transaction success...'});

        await this.raffle.methods.enter().send({
            from: accounts[0],
            value: web3.utils.toWei(this.state.value, 'ether')
        });

        this.setState({message: 'You have been entered!'});
    };

    onClick = async () => {
        const accounts = await web3.eth.getAccounts();

        this.setState({message: 'Waiting on transaction success...'});

        await this.raffle.methods.pickWinner().send({
            from: accounts[0]
        });

        this.setState({message: 'A winner has been picked!'});
    };

    render() {
        return (
            <div>
                <h2>Ringier Raffle</h2>
                <p>
                    This contract is managed by {this.state.manager}.<br/>There are currently{' '}
                    {this.state.players.length} people entered, competing to win{' '}
                    {web3.utils.fromWei(this.state.balance, 'ether')} ETH!
                </p>

                <hr/>

                <form onSubmit={this.onSubmit}>
                    <h4>Want to try your luck?</h4>
                    <div>
                        <label>Amount of ETH to enter </label>
                        <input
                            value={this.state.value}
                            onChange={event => this.setState({value: event.target.value})}
                        />
                    </div>
                    <button>Enter</button>
                </form>

                <hr/>

                <h4>Ready to pick a winner?</h4>
                <button onClick={this.onClick}>Pick a winner!</button>

                <hr/>

                <h1>{this.state.message}</h1>
            </div>
        );
    }

}

export default App;